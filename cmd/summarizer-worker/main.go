package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
	"github.com/meetalodariya/email-thread-summarizer/config"
	emailsummarizer "github.com/meetalodariya/email-thread-summarizer/internal/email_summarizer"
	"github.com/meetalodariya/email-thread-summarizer/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbClient         *gorm.DB
	oauthConfig      *config.OAuth2Config
	sqsClient        *sqs.Client
	queueUrl         string
	openAIKey        string
	openAISummarizer emailsummarizer.EmailSummarizer
)

const USER_EMAIL = "me"

func init() {
	var err error

	env := os.Getenv("ENV")
	if env != "production" {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Failed to load the local .env file.", err)
		}
	}

	dbConfig := config.LoadDBConfig()
	oauthConfig = config.GetRegisterConfig()

	queueUrl = os.Getenv("SQS_QUEUE_URL")
	openAIKey = os.Getenv("OPENAI_API_KEY")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password,
		dbConfig.Database, dbConfig.Port)

	dbClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect the DB.", err)
	}

	dbClient.AutoMigrate(&model.ThreadSummary{})

	openAISummarizer = emailsummarizer.NewOpenAISummarizer(openAIKey)
	err = openAISummarizer.TestAPIKey()
	if err != nil {
		log.Fatal("Failed to connect to Open AI API.", err)
	}
}

func main() {
	// lambda.Start(handler)

	event := events.SQSEvent{Records: []events.SQSMessage{
		{
			MessageId: "123",
			Body:      "{\"userId\":\"6\"}",
		},
	}}
	handler(context.Background(), event)
}
