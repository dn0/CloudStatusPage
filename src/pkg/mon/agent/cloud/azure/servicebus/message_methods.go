package servicebus

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

const (
	messageCount = 5
)

type message struct {
	namespace    string
	queue        string
	topic        string
	subscription string

	client   *azservicebus.Client
	sender   *azservicebus.Sender
	receiver *azservicebus.Receiver
	timeout  time.Duration

	sent     map[string]time.Time
	received map[string]time.Duration
}

func (m *message) String() string {
	if m.queue == "" {
		return fmt.Sprintf("namespace=%s topic=%s subscription=%s", m.namespace, m.topic, m.subscription)
	}
	return fmt.Sprintf("namespace=%s queue=%s", m.namespace, m.queue)
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
	msg := &azservicebus.Message{
		Body:      []byte(mid),
		MessageID: &mid,
	}
	if err := m.sender.SendMessage(ctx, msg, nil); err != nil {
		return fmt.Errorf("message(%s id=%s).Send: %w", m, mid, err)
	}
	m.sent[mid] = now

	return nil
}

func (m *message) receive(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	messages, err := m.receiver.ReceiveMessages(ctx, 1, nil)
	if err != nil {
		return fmt.Errorf("message(%s).Receive: %w", m, err)
	}

	for _, msg := range messages {
		mid := string(msg.Body)
		if err = m.receiver.CompleteMessage(ctx, msg, nil); err != nil {
			return fmt.Errorf("message(%s id=%s).Complete: %w", m, mid, err)
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
