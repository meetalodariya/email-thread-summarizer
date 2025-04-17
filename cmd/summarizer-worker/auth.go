package main

import (
	"context"
	"fmt"
	"log"

	"github.com/meetalodariya/email-thread-summarizer/model"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func authenticateUser(ctx context.Context, userId uint) (*model.User, *oauth2.Token, error) {
	var user model.User
	if result := dbClient.First(&user, userId); result.Error != nil {
		var err error
		if gorm.ErrRecordNotFound == result.Error {
			err = fmt.Errorf("user not found %d: %w", userId, result.Error)
		} else {
			err = fmt.Errorf("failed to fetch user %d: %w", userId, result.Error)
		}

		log.Println(err)
		return nil, nil, err
	}

	oldToken := &oauth2.Token{
		AccessToken:  user.GmailAccessToken,
		RefreshToken: user.GmailRefreshToken,
		Expiry:       user.GmailTokenExpiry,
	}

	token := oldToken
	// Check if access token is already valid
	if !token.Valid() {
		// If not valid then try to refresh the token.
		tokenSource := oauthConfig.Conf.TokenSource(ctx, token)
		newToken, err := tokenSource.Token()

		// If token is failed to refresh then mark the user as unauthenticated in db.
		// Possible reasons: Refresh token expired, user grant revoked, etc.
		// Check GCP consent screen and oauth client config.
		if err != nil {
			log.Println(err)
			log.Println(fmt.Errorf("failed to refresh token for user: %d", userId))

			if dbErr := markUserAsUnauthenticated(userId); dbErr != nil {
				log.Println(fmt.Errorf("failed to update the user as unauthenticated in db: %w", err))
				return nil, nil, dbErr
			}

			return nil, nil, err
		}

		// If token is changed then update the new token in db.
		if newToken.AccessToken != oldToken.AccessToken {
			log.Println(fmt.Printf("Tokens refreshed for the user: %d", userId))

			user.GmailAccessToken = newToken.AccessToken
			user.GmailRefreshToken = newToken.RefreshToken
			user.GmailTokenExpiry = newToken.Expiry

			if err = saveTokenToDb(newToken, userId); err != nil {
				log.Println(err)
				return nil, nil, err
			}

			token = newToken
		}
	}

	return &user, token, nil
}
