package scribe

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/proto"

	"cspage/pkg/db"
	"cspage/pkg/msg"
	"cspage/pkg/pb"
	"cspage/pkg/worker"
)

var errUnknownMessageType = errors.New("unknown message type")

type job struct {
	cfg  *Config
	name string
	dbc  *db.Clients
}

func (j *job) String() string {
	return worker.ConsumerJobPrefix + j.name
}

func (j *job) PreStart(_ context.Context) {}

func (j *job) Start(_ context.Context) {}

func (j *job) Stop(_ context.Context) {}

func (j *job) Shutdown(_ error) {}

func (j *job) Process(ctx context.Context, message *msg.Message) error {
	log := slog.With(
		"job", j.String(),
		"msg_id", message.ID,
		"msg_attrs", message.Attributes,
		"msg_attempt", message.DeliveryAttempt,
	)
	attrs := msg.GetAttrs(message)
	var obj pb.DataModel

	//nolint:exhaustive // Scribe takes care of agent messages only.
	switch attrs.Type {
	case msg.TypeAgent:
		obj = &pb.Agent{}
	case msg.TypePing:
		obj = &pb.Ping{}
	case msg.TypeProbe:
		obj = &pb.Probe{}
	default:
		return fmt.Errorf("%s: %w", attrs.Type, errUnknownMessageType)
	}

	if err := proto.Unmarshal(message.Data, obj); err != nil {
		return fmt.Errorf("proto.Unmarshal(%+v) failed: %w", message.Data, err)
	}

	if err := obj.Save(ctx, j.dbc.Write, attrs); err != nil {
		if db.ErrorIsUniqueViolation(err) {
			log.Info("Duplicate object(s) could not be saved in DB", "obj", obj.Repr(), "err", err)
			return nil // so that the message gets acknowledged
		}
		return fmt.Errorf("obj(%s).Save failed: %w", obj.Repr(), err)
	}
	log.Debug("Object(s) saved in DB", "obj", obj.Repr())
	return nil
}
