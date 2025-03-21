package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	config "github.com/meetalodariya/email-thread-summarizer/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbClient    *gorm.DB
	oauthConfig *config.OAuth2Config
	sqsClient   *sqs.Client
	queueUrl    string
)

func init() {
	var err error
	dbConfig := config.LoadDBConfig()
	oauthConfig = config.GetRegisterConfig()

	queueUrl = os.Getenv("SQS_QUEUE_URL")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=enabled",
		dbConfig.Host, dbConfig.User, dbConfig.Password,
		dbConfig.Database, dbConfig.Port)

	dbClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to db", err)
	}
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event json.RawMessage) error {
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)

	if err != nil {
		err := fmt.Errorf("could not load config: %w", err)
		log.Println(err)

		return err
	}

	sqsClient = sqs.NewFromConfig(awsCfg)

	err = handleRequest(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
