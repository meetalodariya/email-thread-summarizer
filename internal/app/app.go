package app

import (
	"context"
	"log"
	"os"
	"strings"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/meetalodariya/email-thread-summarizer/internal/app/handlers"
	"github.com/meetalodariya/email-thread-summarizer/internal/auth"
	"gorm.io/gorm"
)

type App struct {
	DB  *gorm.DB
	svr *echo.Echo
}

func NewApp(db *gorm.DB) *App {
	return &App{
		DB:  db,
		svr: echo.New(),
	}
}

func (a *App) InitHttpServer() {
	e := a.svr
	handler := handlers.NewHandler(a.DB)

	api := e.Group("/api")

	api.GET("/auth/register/google", handler.HandleRegister)
	api.GET("/auth/register/google/callback", handler.HandleRegisterOAuthCallback)

	api.GET("/auth/login/google", handler.HandleLogin)
	api.GET("/auth/login/google/callback", handler.HandleLoginOAuthCallback)

	protectedG := api.Group("")
	protectedG.Use(echojwt.WithConfig(auth.GetEchoJwtConfig()))

	protectedG.GET("/inbox", handler.HandleGetUserInbox)

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				log.Println("HTTP server shut down")
			} else {
				log.Fatalf("error starting HTTP server: %v", err)
			}
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) error {
	err := a.svr.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
