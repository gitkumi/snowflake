package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSQueue struct {
	client   *sqs.Client
	queueURL string
}

type SQSQueueConfig struct {
	AccessKey string
	SecretKey string
	Region    string
	QueueURL  string
}

func NewSQSQueue(cfg *SQSQueueConfig) *SQSQueue {
	client := createSQSClient(cfg.AccessKey, cfg.SecretKey, cfg.Region)

	return &SQSQueue{
		client:   client,
		queueURL: cfg.QueueURL,
	}
}

func createSQSClient(accessKey, secretKey, region string) *sqs.Client {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))

	cfg := aws.Config{
		Credentials: creds,
		Region:      region,
	}

	return sqs.NewFromConfig(cfg)
}

func (q *SQSQueue) Send(message string) (string, error) {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(q.queueURL),
	}

	result, err := q.client.SendMessage(context.Background(), input)
	if err != nil {
		return "", err
	}

	return *result.MessageId, nil
}

func (q *SQSQueue) Receive(maxMessages int) ([]Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.queueURL),
		MaxNumberOfMessages: int32(maxMessages),
		WaitTimeSeconds:     20, // Long polling
	}

	result, err := q.client.ReceiveMessage(context.Background(), input)
	if err != nil {
		return nil, err
	}

	messages := make([]Message, 0, len(result.Messages))
	for _, msg := range result.Messages {
		messages = append(messages, Message{
			ID:            *msg.MessageId,
			Body:          *msg.Body,
			ReceiptHandle: *msg.ReceiptHandle,
		})
	}

	return messages, nil
}

func (q *SQSQueue) Delete(receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := q.client.DeleteMessage(context.Background(), input)
	return err
}
