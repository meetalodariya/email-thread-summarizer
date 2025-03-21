package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/meetalodariya/email-thread-summarizer/internal/sqs"
	"github.com/meetalodariya/email-thread-summarizer/model"
	"golang.org/x/oauth2"
)

const (
	// In minutes
	SCAN_INTERVAL   = 20
	USER_BATCH_SIZE = 25
)

type UserMessage struct {
	UserID string `json:"userId"`
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	users := make([]model.User, USER_BATCH_SIZE)
	if result := dbClient.Select(
		"id",
		"gmail_access_token",
		"gmail_refresh_token",
		"gmail_token_expiry",
	).Where("last_scanned_timestamp < ?",
		time.Now().Add(-SCAN_INTERVAL*time.Minute),
	).Limit(USER_BATCH_SIZE).Find(&users); result.Error != nil {
		err := fmt.Errorf("failed to fetch users: %w", result.Error)
		log.Println(err)

		return err
	}

	authenticatedUserIds := make([]string, len(users))
	for _, user := range users {
		userId := user.ID
		token := &oauth2.Token{
			AccessToken:  user.GmailAccessToken,
			RefreshToken: user.GmailRefreshToken,
			Expiry:       user.GmailTokenExpiry,
		}

		// Check if access token is already valid
		if !token.Valid() {
			// If not valid then try to refresh the token.
			tokenSource := oauthConfig.Conf.TokenSource(ctx, token)
			newToken, err := tokenSource.Token()

			// If token is failed to refresh then mark the user as unauthenticated in db.
			// Possible reasons: Refresh token expired, user grant revoked, etc.
			// Check GCP consent screen and oauth client config.
			if err != nil {
				log.Println(fmt.Errorf("failed to refresh token for user: %d", userId))

				if err = markUserAsUnauthenticated(userId); err != nil {
					log.Println(fmt.Errorf("failed to update the user as unauthenticated in db: %w", err))
					return err
				}

				log.Println(fmt.Printf("Tokens refreshed for the user: %d", userId))

				continue
			}

			// If token is changed then update the new token in db.
			if newToken.AccessToken != token.AccessToken {
				if err = saveTokenToDb(newToken, userId); err != nil {
					log.Println(err)
					return err
				}
			}
		}

		userMessage := UserMessage{
			UserID: fmt.Sprintf("%d", userId),
		}

		jsonData, err := json.Marshal(userMessage)
		if err != nil {
			log.Println(err)
			return err
		}

		authenticatedUserIds = append(authenticatedUserIds, string(jsonData))
	}

	config := sqs.Config{
		QueueURL: queueUrl,
		// using the default batch size and max retries
	}

	c := sqs.NewMessageQueueClient(sqsClient, config)

	if err := c.SendMessages(ctx, authenticatedUserIds); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
