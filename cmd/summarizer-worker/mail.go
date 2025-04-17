package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"sort"
	"time"

	"github.com/meetalodariya/email-thread-summarizer/model"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Recipient struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Email represents a simplified email structure
type Email struct {
	ID      string
	Subject string

	Date     time.Time
	Body     string
	ThreadID string
	From     Recipient
	To       Recipient
}

// processEmails fetches and processes new emails.
func processEmails(ctx context.Context, usr *model.User, token *oauth2.Token) error {
	var err error
	srv, err := gmail.NewService(ctx, option.WithTokenSource(oauthConfig.Conf.TokenSource(ctx, token)))
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %w", err)
	}

	log.Printf("Processing the emails posted between %s and %s of the user: %s",
		getTimeString(usr.LastScannedTimestamp),
		getTimeString(time.Now()),
		usr.Email,
	)

	// Fetch emails
	emails, err := fetchNewEmails(ctx, srv, usr.LastScannedTimestamp)
	if err != nil {
		log.Printf("An error occurred while fetching emails: %v\n", err)

		return err
	}

	if len(emails) == 0 {
		return nil
	}

	for _, email := range emails {
		if err := summarizeEmail(ctx, email, usr); err != nil {
			log.Printf("couldn't process email: %s", email.ID)
			return err
		}
	}

	lastEmail := emails[len(emails)-1]
	if err := updateUserLastProcessed(usr, lastEmail); err != nil {
		return fmt.Errorf("failed to update user's last processed info: %w", err)
	}

	return nil
}

func fetchNewEmails(ctx context.Context, srv *gmail.Service, since time.Time) ([]*Email, error) {
	query := fmt.Sprintf("after:%d", since.Unix())
	var emails []*Email
	err := srv.Users.Messages.List(USER_EMAIL).Q(query).Pages(ctx,
		func(res *gmail.ListMessagesResponse) error {
			for _, msgData := range res.Messages {
				msg, err := srv.Users.Messages.Get(USER_EMAIL, msgData.Id).Do()
				if err != nil {
					return err
				}
				emails = append(emails, extractEmailData(msg))
			}
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	// Sort emails by timestamp (descending)
	sort.Slice(emails, func(i, j int) bool {
		return emails[i].Date.After(emails[j].Date)
	})
	return emails, nil
}

func extractEmailData(msg *gmail.Message) *Email {
	email := &Email{
		ID: msg.Id,
	}

	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "Subject":
			email.Subject = header.Value
		case "From":
			addr, err := mail.ParseAddress(header.Value)
			if err == nil {
				email.From = Recipient{
					Email: addr.Address,
					Name:  addr.Name,
				}
			}
		case "To":
			addr, err := mail.ParseAddress(header.Value)
			if err == nil {
				email.To = Recipient{
					Email: addr.Address,
					Name:  addr.Name,
				}
			}
		}
	}

	email.Body = parseMessageBody(msg.Payload)
	email.Date = time.UnixMilli(msg.InternalDate)
	email.ThreadID = msg.ThreadId

	return email
}

func summarizeEmail(ctx context.Context, e *Email, usr *model.User) error {
	log.Printf("Processing email: %s - Subject: %s", e.ID, e.Subject)
	ts, err := getOrCreateThreadSummary(e.ThreadID)
	if err != nil {
		log.Printf("Failed to get or create thread summary: %v", err)
		return err
	}

	// Skip if email is already processed
	for _, processedId := range ts.ProcessedEmailIds {
		if processedId == e.ID {
			log.Printf("Email %s already processed, skipping", e.ID)
			return nil
		}
	}

	if ts.GmailThreadId != "" {
		ts.Summary, err = openAISummarizer.SummarizeEmailIncremental(ctx, ts.Summary, e.Body, e.Subject, e.From.Name)
		if err != nil {
			log.Printf("Failed to summarize email incrementally: %v", err)
			return err
		}
	} else {
		ts.Summary, err = openAISummarizer.SummarizeEmail(ctx, e.Body, e.Subject, e.From.Name)
		if err != nil {
			log.Printf("Failed to summarize email: %v", err)
			return err
		}
		ts.GmailThreadId = e.ThreadID
		ts.UserID = usr.ID
	}

	ts.ThreadSubject = e.Subject
	ts.MostRecentEmailTimestamp = e.Date
	ts.ProcessedEmailIds = append(ts.ProcessedEmailIds, e.ID)
	ts.UpdatedAt = e.Date

	if e.From.Email != usr.Email && e.From.Name != "" {
		ts.Recipients = addUnique(ts.Recipients, e.From.Name)
	}

	if e.To.Email != usr.Email && e.To.Name != "" {
		ts.Recipients = addUnique(ts.Recipients, e.To.Name)
	}

	if err := saveThreadSummary(ts); err != nil {
		log.Printf("Failed to save thread summary: %v", err)
		return err
	}

	return nil
}

func parseMessageBody(payload *gmail.MessagePart) string {
	if payload.Body != nil && payload.Body.Data != "" {
		// Decode base64 URL-encoded body
		decoded, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			log.Printf("Error decoding body: %v", err)
			return ""
		}
		return string(decoded)
	}

	// If it's multipart, find the plain text or HTML body
	for _, part := range payload.Parts {
		if part.MimeType == "text/plain" || part.MimeType == "text/html" {
			return parseMessageBody(part)
		}
	}
	return ""
}

func updateUserLastProcessed(usr *model.User, email *Email) error {
	usr.LastScannedTimestamp = email.Date
	usr.LastProcessedMail = email.ID

	if err := dbClient.Save(usr).Error; err != nil {
		return fmt.Errorf("failed to update user's last scanned timestamp: %w", err)
	}
	return nil
}
