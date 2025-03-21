package sqs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// MessageQueueClient represents a client for sending messages to SQS with configuration
type MessageQueueClient struct {
	client     *sqs.Client
	queueURL   string
	batchSize  int
	maxRetries int
	retryDelay time.Duration
}

// Config holds the configuration for MessageQueueClient
type Config struct {
	QueueURL   string
	BatchSize  int
	MaxRetries int
	RetryDelay time.Duration
}

// NewMessageQueueClient creates a new message queue client with the given configuration
func NewMessageQueueClient(client *sqs.Client, config Config) *MessageQueueClient {
	if config.BatchSize <= 0 {
		config.BatchSize = 10 // default batch size
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3 // default max retries
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = time.Second * 2 // default retry delay
	}

	return &MessageQueueClient{
		client:     client,
		queueURL:   config.QueueURL,
		batchSize:  config.BatchSize,
		maxRetries: config.MaxRetries,
		retryDelay: config.RetryDelay,
	}
}

// SendMessages sends multiple messages to SQS in batches
func (s *MessageQueueClient) SendMessages(ctx context.Context, messages []string) error {
	for i := 0; i < len(messages); i += s.batchSize {
		end := i + s.batchSize
		if end > len(messages) {
			end = len(messages)
		}

		// Send batch messages with retry
		if err := s.sendMessageBatchWithRetry(ctx, messages[i:end]); err != nil {
			return fmt.Errorf("failed to send message batch: %w", err)
		}
	}

	return nil
}

// sendMessageBatchWithRetry attempts to send a batch of messages with retries
func (s *MessageQueueClient) sendMessageBatchWithRetry(ctx context.Context, messages []string) error {
	var lastErr error
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if err := s.sendMessageBatch(ctx, messages); err != nil {
			lastErr = err
			log.Printf("Attempt %d failed: %v", attempt+1, err)

			if attempt < s.maxRetries {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(s.retryDelay):
					continue
				}
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("failed after %d attempts, last error: %w", s.maxRetries, lastErr)
}

// sendMessageBatch sends a batch of messages to SQS
func (s *MessageQueueClient) sendMessageBatch(ctx context.Context, msgBodies []string) error {
	entries := make([]types.SendMessageBatchRequestEntry, len(msgBodies))

	for i, msgBody := range msgBodies {

		entries[i] = types.SendMessageBatchRequestEntry{
			Id:          aws.String(fmt.Sprintf("%d", i)),
			MessageBody: aws.String(msgBody),
		}
	}

	result, err := s.client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(s.queueURL),
		Entries:  entries,
	})

	if err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	// Handle failed messages
	if len(result.Failed) > 0 {
		var failedIds []string
		for _, failed := range result.Failed {
			failedIds = append(failedIds, *failed.Id)
		}
		return fmt.Errorf("some messages failed to send: %v", failedIds)
	}

	return nil
}
