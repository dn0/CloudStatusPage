package pubsub

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	cloudPubsub "cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

const (
	pubsubEndpoint = "%s-pubsub.googleapis.com:443"
)

var errReceivedAllMessages = errors.New("received all messages")

type message struct {
	count          int
	projectID      string
	region         string
	topicID        string
	subscriptionID string

	client       *cloudPubsub.Client
	topic        *cloudPubsub.Topic
	subscription *cloudPubsub.Subscription
	timeout      time.Duration

	mu       sync.Mutex
	sent     map[string]time.Time
	received map[string]time.Duration
}

func (m *message) String() string {
	return fmt.Sprintf("topic=%s subscription=%s", m.topicID, m.subscriptionID)
}

func (m *message) sentAll() bool {
	return m.count == len(m.sent)
}

func (m *message) receivedAll() bool {
	return m.count == len(m.received)
}

func (m *message) new() {
	m.sent = make(map[string]time.Time, m.count)
	m.received = make(map[string]time.Duration, m.count)
}

func (m *message) init(ctx context.Context) error {
	var err error
	regionalEndpoint := fmt.Sprintf(pubsubEndpoint, m.region)
	m.client, err = cloudPubsub.NewClient(ctx, m.projectID, option.WithEndpoint(regionalEndpoint))
	if err != nil {
		return fmt.Errorf("message(%s endpoint=%s).NewClient: %w", m, regionalEndpoint, err)
	}

	m.topic = m.client.Topic(m.topicID)
	m.topic.PublishSettings = cloudPubsub.PublishSettings{
		CountThreshold: 1,
	}

	m.subscription = m.client.Subscription(m.subscriptionID)
	m.subscription.ReceiveSettings = cloudPubsub.ReceiveSettings{
		MaxOutstandingMessages: 1,
	}

	return nil
}

func (m *message) close() {
	if m.topic != nil {
		m.topic.Stop()
	}
	if m.client != nil {
		_ = m.client.Close()
		m.client = nil
		m.topic = nil
		m.subscription = nil
	}
}

func (m *message) send(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	now := time.Now()
	mid := strconv.FormatInt(now.UnixNano(), 10)
	res := m.topic.Publish(ctx, &cloudPubsub.Message{Data: []byte(mid)})
	if _, err := res.Get(ctx); err != nil {
		return fmt.Errorf("message(%s id=%s).Publish: %w", m, mid, err)
	}
	m.sent[mid] = now

	return nil
}

func (m *message) receive(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	ctx, stop := context.WithCancelCause(ctx)
	defer stop(nil)

	err := m.subscription.Receive(ctx, func(_ context.Context, msg *cloudPubsub.Message) {
		msg.Ack()
		mid := string(msg.Data)
		m.mu.Lock() // Yes, this makes the whole thing slow.
		defer m.mu.Unlock()
		sendTime, exists := m.sent[mid]
		if exists {
			m.received[mid] = time.Since(sendTime)
			delete(m.sent, mid)
		}
		if m.receivedAll() {
			stop(errReceivedAllMessages)
		}
	})
	if err != nil && !errors.Is(context.Cause(ctx), errReceivedAllMessages) {
		return fmt.Errorf("message(%s).Receive: %w", m, err)
	}

	return nil
}
