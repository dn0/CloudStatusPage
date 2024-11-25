package sqs

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	messageCount              = 5
	receiveWaitTimeoutSeconds = 5
)

type message struct {
	queueURL string
	client   *sqs.Client
	timeout  time.Duration

	sent     map[string]time.Time
	received map[string]time.Duration
}

func (m *message) String() string {
	return "url=" + m.queueURL
}

func (m *message) sentAll() bool {
	return messageCount == len(m.sent)
}

func (m *message) receivedAll() bool {
	return messageCount == len(m.received)
}

func (m *message) new() {
	m.sent = make(map[string]time.Time, messageCount)
	m.received = make(map[string]time.Duration, messageCount)
}

func (m *message) send(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	now := time.Now()
	mid := strconv.FormatInt(now.UnixNano(), 10)
	if _, err := m.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(m.queueURL),
		MessageBody: aws.String(mid),
	}); err != nil {
		return fmt.Errorf("message(%s id=%s).Send: %w", m, mid, err)
	}
	m.sent[mid] = now

	return nil
}

func (m *message) receive(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	res, err := m.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(m.queueURL),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     receiveWaitTimeoutSeconds,
	})
	if err != nil {
		return fmt.Errorf("message(%s).Receive: %w", m, err)
	}

	for _, msg := range res.Messages {
		mid := *msg.Body
		if _, err = m.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(m.queueURL),
			ReceiptHandle: msg.ReceiptHandle,
		}); err != nil {
			return fmt.Errorf("message(%s id=%s).Delete: %w", m, mid, err)
		}
		sendTime, exists := m.sent[mid]
		if !exists {
			continue
		}
		m.received[mid] = time.Since(sendTime)
		delete(m.sent, mid)
	}

	return nil
}
