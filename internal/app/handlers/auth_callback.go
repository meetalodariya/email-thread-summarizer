package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/meetalodariya/email-thread-summarizer/config"
	"github.com/meetalodariya/email-thread-summarizer/internal/auth"
	"github.com/meetalodariya/email-thread-summarizer/model"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type UserInfo struct {
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
}

type googleAuthenticationResponse struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type googleAuthenticationError struct {
	Message string `json:"message"`
}

func (h *Handler) HandleGoogleAuthentication(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	code, err := extractAuthCode(eCtx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	oauthConfig := config.GetRegisterConfig()
	tok, err := oauthConfig.Conf.Exchange(ctx, code)
	if err != nil {
		return authError(http.StatusUnauthorized, eCtx)
	}

	userInfo, err := fetchUserInfo(ctx, oauthConfig, tok)
	if err != nil {
		return authError(http.StatusUnauthorized, eCtx)
	}

	user, err := fetchUser(h.DB, userInfo.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return authError(http.StatusInternalServerError, eCtx)
	}

	if user != nil {
		return respondWithToken(user.FirstName+" "+user.LastName, user.ID, eCtx)
	}

	newUser := createUserFromInfo(userInfo, tok)
	if err := h.DB.Create(&newUser).Error; err != nil {
		return authError(http.StatusInternalServerError, eCtx)
	}

	return respondWithToken(newUser.FirstName+" "+newUser.LastName, newUser.ID, eCtx)
}

func extractAuthCode(eCtx echo.Context) (string, error) {
	var body map[string]string
	if err := json.NewDecoder(eCtx.Request().Body).Decode(&body); err != nil {
		return "", err
	}
	return body["code"], nil
}

func createUserFromInfo(info *UserInfo, tok *oauth2.Token) model.User {
	return model.User{
		FirstName:            info.FirstName,
		LastName:             info.LastName,
		Picture:              info.Picture,
		Email:                info.Email,
		GmailRefreshToken:    tok.RefreshToken,
		GmailAccessToken:     tok.AccessToken,
		GmailTokenExpiry:     tok.Expiry,
		IsGmailTokenValid:    true,
		LastScannedTimestamp: time.Now().AddDate(0, 0, -1),
	}
}

func respondWithToken(name string, userID uint, eCtx echo.Context) error {
	token, err := auth.GenerateToken(userID)
	if err != nil {
		return authError(http.StatusInternalServerError, eCtx)
	}
	return eCtx.JSON(http.StatusOK, googleAuthenticationResponse{Name: name, Token: token})
}

func authError(status int, eCtx echo.Context) error {
	return eCtx.JSON(status, googleAuthenticationError{Message: "Something went wrong!"})
}

func fetchUser(db *gorm.DB, email string) (*model.User, error) {
	var user model.User
	result := db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func fetchUserInfo(ctx context.Context, config *config.OAuth2Config, t *oauth2.Token) (*UserInfo, error) {
	client := config.Conf.Client(ctx, t)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}
