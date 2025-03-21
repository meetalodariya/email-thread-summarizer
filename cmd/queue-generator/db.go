package main

import (
	"fmt"
	"log"
	"time"

	"github.com/meetalodariya/email-thread-summarizer/model"
	"golang.org/x/oauth2"
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

	if result := dbClient.Model(&user).Updates(model.User{
		IsGmailTokenValid: false,
		GmailAccessToken:  "",
		GmailRefreshToken: "",
		GmailTokenExpiry:  time.Time{},
	}); result.Error != nil {
		err := fmt.Errorf("failed to save new token: %w", result.Error)
		log.Println(err)
		return err
	}

	return nil
}
