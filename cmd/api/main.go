package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/meetalodariya/email-thread-summarizer/config"
	"github.com/meetalodariya/email-thread-summarizer/internal/app"
	"github.com/meetalodariya/email-thread-summarizer/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!", err)
		os.Exit(0)
	}

	dbConfig := config.LoadDBConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password,
		dbConfig.Database, dbConfig.Port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	db.AutoMigrate(&model.User{})

	if err != nil {
		log.Fatal("Could not connect to db", err)
	}
}

func main() {
	app := app.NewApp(db)

	// Start http server.
	app.InitHttpServer()

	quit := make(chan os.Signal, 1)
	// Listen for SIGINT (Ctrl+C) and SIGTERM
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until signal is received
	<-quit
	log.Println("Shutting down server...")

	shutdown(app)

	log.Println("Server exited gracefully")
}

func shutdown(app *app.App) {
	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown app
	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
