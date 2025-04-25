package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth2Config struct {
	Conf *oauth2.Config
}

func GetRegisterScopes() []string {
	return []string{
		"https://mail.google.com",
		"https://www.googleapis.com/auth/gmail.readonly",
		"https://www.googleapis.com/auth/gmail.modify",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
}

func GetRegisterConfig() *OAuth2Config {
	return &OAuth2Config{
		Conf: &oauth2.Config{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			RedirectURL:  os.Getenv("CLIENT_REGISTER_CALLBACK_URL"),
			Scopes:       GetRegisterScopes(),
			Endpoint:     google.Endpoint,
		},
	}
}
