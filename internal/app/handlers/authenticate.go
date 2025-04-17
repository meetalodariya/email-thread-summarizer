package handlers

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/meetalodariya/email-thread-summarizer/config"
	"golang.org/x/oauth2"
)

func (h *Handler) HandleGoogleAuthenticationInit(eCtx echo.Context) error {
	req := eCtx.Request()

	oauthConfig := config.GetRegisterConfig()

	authURL := oauthConfig.Conf.AuthCodeURL("example", oauth2.AccessTypeOffline)

	u, err := url.Parse(authURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not generate auth URL.")
	}

	q := u.Query()
	q.Add("prompt", "consent")

	u.RawQuery = q.Encode()

	authURL = u.String()

	http.Redirect(eCtx.Response().Writer, req, authURL, http.StatusFound)

	return nil
}
