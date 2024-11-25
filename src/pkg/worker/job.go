package worker

import (
	"context"
	"time"

	"cspage/pkg/msg"
)

type job interface {
	String() string
	PreStart(context.Context)
	Start(context.Context)
	Stop(context.Context)
	Shutdown(error)
}

type TickerJob interface {
	job
	Do(context.Context, Tick)
}

type ConsumerJob interface {
	job
	Process(context.Context, *msg.Message) error
}

type ReceiverJob[T any] interface {
	job
	Enabled() bool
	Process(context.Context, T) error
}

type Tick struct {
	ID    string        `json:"id"`
	Time  time.Time     `json:"time"`
	Drift time.Duration `json:"drift"`
	Took  time.Duration `json:"took"`
}
