package sqs

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	messageProbeSleep         = 10 * time.Millisecond // Sleep between send and receive. Subtracted from the final message latency.
	messageProbeName          = "aws_sqs_message"     // doc="Amazon SQS Message"
	messageProbeActionSend    = 60                    // name="sqs.message.send" doc="Sends a couple of small messages to a queue" url="https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_SendMessage.html"
	messageProbeActionReceive = 70                    // name="sqs.message.receive" doc="Receives all messages sent by `sqs.message.send` and removes them from queue" url="https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ReceiveMessage.html"
	messageProbeActionLatency = 80                    // name="sqs.message.latency" doc="Average delay between sending and receiving of one message" url=""
)

type MessageProbe[T agent.AWS] struct {
	cfg      *agent.AWSConfig
	message  *message
	warmedUp bool
}

func NewMessageProbe[T agent.AWS](cfg *agent.AWSConfig, awsConfig *aws.Config) *MessageProbe[T] {
	return &MessageProbe[T]{
		cfg: cfg,
		message: &message{
			client:  sqs.NewFromConfig(*awsConfig),
			timeout: cfg.ProbeTimeout,
		},
	}
}

func (p *MessageProbe[T]) String() string {
	return messageProbeName
}

func (p *MessageProbe[T]) Start(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, p.cfg.ProbeTimeout)
	defer cancel()

	if res, err := p.message.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(p.cfg.Cloud.SQSQueueName),
	}); err == nil && res.QueueUrl != nil {
		p.message.queueURL = *res.QueueUrl
	} else {
		agent.DieLog(
			p.log(),
			"Could not fetch AWS SQS queue URL",
			"queue_name", p.cfg.Cloud.SQSQueueName,
			"api_response", res,
			"err", err,
		)
	}

	p.log().Info("Probe initialized")
}

func (p *MessageProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(messageProbeActionSend),
		pb.NewResult(messageProbeActionReceive),
		pb.NewResult(messageProbeActionLatency),
	}
	p.message.new()

	res[0].Store(p.sendMessages(ctx))
	if res[0].Failed() {
		return res
	}

	time.Sleep(messageProbeSleep)

	res[1].Store(p.receiveMessages(ctx))

	if p.message.receivedAll() {
		res[2].Store(p.getLatency(), nil)
	}

	if !p.warmedUp {
		p.log().Debug("Cold start => discarding all results")
		p.warmedUp = true
		return nil
	}

	return res
}

func (p *MessageProbe[T]) Stop(_ context.Context) {}

//nolint:dupl // We can tolerate the duplicate code between receiveMessages and sendMessages.
func (p *MessageProbe[T]) sendMessages(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Sending messages to SQS...")
	retAll := pb.EmptyResultTime()
	var err error
	var num int

	for !p.message.sentAll() {
		var ret pb.ResultTime
		ret, err = pb.Timeit(p.message.send, ctx)
		num++
		//goland:noinspection GoDfaErrorMayBeNotNil
		retAll.Took += ret.Took
		if err == nil {
			p.log().Debug("Sent message to SQS", "num", num, "took", ret.Took)
		} else {
			p.log().Error("Could not send message to SQS", "num", num, "took", ret.Took, "err", err)
			break
		}
	}

	retAll.Took /= time.Duration(num)
	p.log().Debug("Sent messages to SQS", "sent", num, "took", retAll.Took)
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return retAll, err
}

//nolint:dupl // We can tolerate the duplicate code between receiveMessages and sendMessages.
func (p *MessageProbe[T]) receiveMessages(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Receiving messages from SQS...")
	retAll := pb.EmptyResultTime()
	var err error
	var num int

	for !p.message.receivedAll() {
		var ret pb.ResultTime
		ret, err = pb.Timeit(p.message.receive, ctx)
		num++
		//goland:noinspection GoDfaErrorMayBeNotNil
		retAll.Took += ret.Took
		if err == nil {
			p.log().Debug("Received message from SQS", "num", num, "took", ret.Took)
		} else {
			p.log().Error("Could not receive message from SQS", "num", num, "took", ret.Took, "err", err)
			break
		}
	}

	retAll.Took /= time.Duration(num)
	p.log().Debug("Received messages from SQS", "received", num, "took", retAll.Took)
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return retAll, err
}

func (p *MessageProbe[T]) getLatency() pb.ResultTime {
	ret := pb.EmptyResultTime()
	if len(p.message.received) == 0 {
		return ret
	}

	var sum time.Duration
	for _, took := range p.message.received {
		sum += took - messageProbeSleep
	}

	ret.Took = sum / time.Duration(len(p.message.received))
	p.log().Debug("Average message send/receive latency", "took", ret.Took)
	return ret
}

func (p *MessageProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "sqs", p.message.String(), "count", messageCount)
}
