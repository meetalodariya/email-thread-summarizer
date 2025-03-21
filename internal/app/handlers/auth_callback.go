package handlers

import (
	"context"
	"encoding/json"
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

type oauthLoginCallbackResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func (h *Handler) HandleRegisterOAuthCallback(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	req := eCtx.Request()
	q := req.URL.Query()

	receivedState := q.Get("state")
	if receivedState != "example" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid state")
	}

	oauthConfig := config.GetRegisterConfig()

	code := q.Get("code")
	tok, err := oauthConfig.Conf.Exchange(ctx, code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Token exchange failed")
	}

	user, err := fetchUserInfo(ctx, oauthConfig, tok)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not fetch user data")
	}

	userExists, err := userEmailExists(h.DB, user.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	if userExists {
		return echo.NewHTTPError(http.StatusConflict, "User already exists.")
	}

	if result := h.DB.Create(&model.User{
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Picture:              user.Picture,
		Email:                user.Email,
		GmailRefreshToken:    tok.RefreshToken,
		GmailAccessToken:     tok.AccessToken,
		GmailTokenExpiry:     tok.Expiry,
		IsGmailTokenValid:    true,
		LastScannedTimestamp: time.Now().AddDate(0, -1, 0),
	}); result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not save the user data")
	}

	return eCtx.JSON(http.StatusCreated, JsonResponse{Data: "User successfully authenticated."})
}

func (h *Handler) HandleLoginOAuthCallback(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	req := eCtx.Request()
	q := req.URL.Query()

	receivedState := q.Get("state")
	if receivedState != "example" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid state")
	}

	oauthConfig := config.GetLoginConfig()

	code := q.Get("code")
	tok, err := oauthConfig.Conf.Exchange(ctx, code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Token exchange failed")
	}

	user, err := fetchUserInfo(ctx, oauthConfig, tok)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not fetch user data")
	}
	println(user.Email)

	usr, err := fetchUser(h.DB, user.Email)
	if err == gorm.ErrRecordNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "User not found.")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	token, err := auth.GenerateToken(usr.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not generate token")
	}

	return eCtx.JSON(http.StatusOK, oauthLoginCallbackResponse{Message: "User successfully authenticated.", Token: token})
}

func userEmailExists(db *gorm.DB, email string) (bool, error) {
	_, err := fetchUser(db, email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
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
