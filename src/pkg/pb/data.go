package pb

import (
	"context"

	"google.golang.org/protobuf/proto"

	"cspage/pkg/db"
	"cspage/pkg/msg"
)

type DataModel interface {
	proto.Message

	ID() string
	Repr() string
	Save(context.Context, db.Client, *msg.Attrs) error
}
