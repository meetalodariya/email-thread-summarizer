package app

import (
	"context"
	"log"
	"os"
	"strings"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

	e.Use(middleware.CORS())
	handler := handlers.NewHandler(a.DB)

	api := e.Group("/api")

	api.GET("/auth/google", handler.HandleGoogleAuthenticationInit)
	api.POST("/auth/google", handler.HandleGoogleAuthentication)

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
