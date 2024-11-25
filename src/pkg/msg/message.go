package msg

import (
	"cloud.google.com/go/pubsub"
	"google.golang.org/protobuf/proto"

	"cspage/pkg/config"
)

const (
	attrMsgType = "msg_type"
	attrCloud   = "cloud"
	attrRegion  = "region"

	TypeAgent    Type = "AGENT"
	TypePing     Type = "PING"
	TypeProbe    Type = "PROBE"
	TypeAlert    Type = "ALERT"
	TypeIncident Type = "INCIDENT"
)

type Type string

type Message = pubsub.Message

type Attrs struct {
	Type   Type
	Cloud  string
	Region string
}

func (a *Attrs) toProto() map[string]string {
	return map[string]string{
		attrMsgType: string(a.Type),
		attrCloud:   a.Cloud,
		attrRegion:  a.Region,
	}
}

func GetAttrs(m *Message) *Attrs {
	return &Attrs{
		Type:   Type(m.Attributes[attrMsgType]),
		Cloud:  m.Attributes[attrCloud],
		Region: m.Attributes[attrRegion],
	}
}

func NewAttrs(msgType Type, cloud, region string) *Attrs {
	return &Attrs{
		Type:   msgType,
		Cloud:  cloud,
		Region: region,
	}
}

func NewMessage(pbMessage proto.Message, attrs *Attrs) *Message {
	return &Message{
		Attributes: attrs.toProto(),
		Data:       pbMessageToData(pbMessage, attrs.Type),
	}
}

func pbMessageToData(pbMessage proto.Message, msgType Type) []byte {
	data, err := proto.Marshal(pbMessage)
	if err != nil {
		config.Die("Could not marshal protobuf message", "msg_type", msgType, "err", err)
	}
	return data
}
