package emailsummarizer

import (
	"context"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

const OPENAI_MODEL = openai.GPT3Dot5Turbo

// EmailSummarizer defines the behavior for summarizing emails.
type EmailSummarizer interface {
	SummarizeEmail(ctx context.Context, emailBody string, subject string, from string) (string, error)
	SummarizeEmailIncremental(ctx context.Context, oldSummary string, emailBody string, subject string, from string) (string, error)
	TestAPIKey() error
}

// OpenAISummarizer implements EmailSummarizer using OpenAI's API.
type OpenAISummarizer struct {
	client *openai.Client
}

// NewOpenAISummarizer creates a new instance of OpenAISummarizer.
func NewOpenAISummarizer(apiKey string) *OpenAISummarizer {
	return &OpenAISummarizer{
		client: openai.NewClient(apiKey),
	}
}

func (o *OpenAISummarizer) TestAPIKey() error {
	_, err := o.client.ListModels(context.Background())
	if err != nil {
		log.Fatal("Failed to connect to Open AI API.", err)
		return err
	}
	return nil
}

func (o *OpenAISummarizer) SummarizeEmail(ctx context.Context, emailBody string, subject string, from string) (string, error) {
	userPrompt := fmt.Sprintf(`
	  Subject: %s,
	  From: %s

	  Body: 
	  %s
	`, subject, from, emailBody)

	chatCompletion, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: OPENAI_MODEL,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `
Summarize the following email in a concise and professional manner. 
Use Markdown in response

Instructions:
1. Respond in markdown syntax as per the 'Format' section provided below. 

2. Highlight important dates/times, locations or any important info in markdown syntax.

3. Provide 'Urgency Score' for the email out of 10. (replace the score with 'score' placeholder in format.)

4. Provide the detailed summary of the email under 'Summary of the thread' in bullet points. (replace the score with 'summary' placeholder in format.)

5. Provide section for the 'Action Items' for action items/next steps outlined in the email. (replace the score with 'actionItems' placeholder in format.)

6. Do not include long links. 

Format:
## Urgency Score: {score}/10 

## Summary of the thread: 
{summary}

## Action Items: 
{actionItems}

Provide separate section for the 'action items'. Gmail Summarizer (gsummarizer@gmail.com) is the user of this API.
Ignore the old conversations that start with '>'. 
					`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to summarize email: %w", err)
	}
	content := chatCompletion.Choices[0].Message.Content

	if content == "" {
		return "", fmt.Errorf("no response from OpenAI")

	}

	return content, nil
}

func (o *OpenAISummarizer) SummarizeEmailIncremental(
	ctx context.Context,
	oldSummary string,
	emailBody string,
	subject string,
	from string) (string, error) {
	userPrompt := fmt.Sprintf(`
Current Thread Summary: 
%s		
---------------------------------------------------------
New Email:
Subject: %s,
From: %s
Body: 
%s
	`, oldSummary, subject, from, emailBody)

	chatCompletion, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: OPENAI_MODEL,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `
Given Current Thread Summary of email thread, summarize new email on the existing thread and provide the overall summary of the thread in the format 
given below. Use Markdown in response.

Instructions:
1. Respond in markdown syntax as per the 'Format' section provided below. 

2. Highlight important dates/times, locations or any important info in markdown syntax.

3. Provide 'Urgency Score' for the email out of 10. (replace the score with 'score' placeholder in format.)

4. Provide the detailed updated summary of the email thread under 'Summary of the thread' in bullet points. (replace the score with 'summary' placeholder in format.)

5. Provide section for the 'Action Items' for updated action items/next steps outlined in the email in bullet points. (replace the score with 'actionItems' placeholder in format.)

6. Do not include long links. 

-------

Format:
## Urgency Score: {score}/10 

## Summary of the thread: 
{summary}

## Action Items: 
{actionItems}

Gmail Summarizer (gsummarizer@gmail.com) is the user of this API.
Ignore the old conversations in the new email that start with '>'.`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to summarize email: %w", err)
	}

	if chatCompletion.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("no response from OpenAI")

	}

	return chatCompletion.Choices[0].Message.Content, nil
}
