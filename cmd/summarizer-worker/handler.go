package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var err error
	// awsCfg, err := config.LoadDefaultConfig(ctx)

	// if err != nil {
	// 	err := fmt.Errorf("could not load config: %w", err)
	// 	log.Println(err)

	// 	return err
	// }

	// sqsClient = sqs.NewFromConfig(awsCfg)

	message := sqsEvent.Records[0]
	err = processMessage(ctx, message)
	if err != nil {
		return err
	}

	// Delete message after successful processing
	// _, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
	// 	QueueUrl:      aws.String(queueUrl),
	// 	ReceiptHandle: aws.String(message.ReceiptHandle),
	// })

	// if err != nil {
	// 	log.Printf("Failed to delete message: %v", err)
	// 	return err
	// }

	log.Printf("Successfully processed and deleted message: %s", message.MessageId)

	return nil
}

type Payload struct {
	UserId string `json:"userId"`
}

func processMessage(ctx context.Context, message events.SQSMessage) error {
	log.Printf("Processing message: %s", message.Body)

	var payload Payload
	err := json.Unmarshal([]byte(message.Body), &payload)
	if err != nil {
		log.Printf("Failed to parse the message: %v", err)
		return err
	}

	userId, err := strconv.Atoi(payload.UserId)
	if err != nil {
		log.Printf("Failed to parse the userId: %v", err)
		return err
	}

	usr, token, err := authenticateUser(ctx, uint(userId))
	if err != nil {
		return err
	}

	processEmails(ctx, usr, token)

	return nil
}
