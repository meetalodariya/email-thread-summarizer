package main

import (
	"fmt"
	"log"

	"github.com/meetalodariya/email-thread-summarizer/model"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func saveTokenToDb(token *oauth2.Token, userId uint) error {
	var user model.User
	user.ID = userId

	if result := dbClient.Model(&user).Updates(model.User{
		GmailAccessToken:  token.AccessToken,
		GmailRefreshToken: token.RefreshToken,
		GmailTokenExpiry:  token.Expiry,
	}); result.Error != nil {
		err := fmt.Errorf("failed to save new token: %w", result.Error)
		log.Println(err)
		return err
	}

	return nil
}

func markUserAsUnauthenticated(userId uint) error {
	var user model.User
	user.ID = userId

	if result := dbClient.Model(&user).Updates(map[string]any{
		"gmail_access_token":   "",
		"gmail_refresh_token":  "",
		"gmail_token_expiry":   "",
		"is_gmail_token_valid": false,
	}); result.Error != nil {
		err := fmt.Errorf("failed to save new token: %w", result.Error)
		log.Println(err)
		return err
	}

	return nil
}

// saveThreadSummary saves or updates a thread summary in the database.
func saveThreadSummary(ts *model.ThreadSummary) error {
	if ts.GmailThreadId != "" {
		result := dbClient.Save(ts)
		return result.Error
	}

	result := dbClient.Create(ts)
	return result.Error
}

// getOrCreateThreadSummary retrieves an existing thread summary or creates a new one.
func getOrCreateThreadSummary(threadID string) (*model.ThreadSummary, error) {
	ts := &model.ThreadSummary{}
	result := dbClient.Where("gmail_thread_id = ?", threadID).First(ts)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &model.ThreadSummary{}, nil
		}
		return nil, result.Error
	}
	return ts, nil
}
