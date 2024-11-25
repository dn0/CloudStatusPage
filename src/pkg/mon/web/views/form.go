package views

import (
	"github.com/go-playground/form/v4"
)

const (
	maxArraySize uint = 10
)

type FormDecoder = form.Decoder

type FormEncoder = form.Encoder

func NewFormDecoder() *FormDecoder {
	decoder := form.NewDecoder()
	decoder.SetMaxArraySize(maxArraySize)
	return decoder
}

func NewFormEncoder() *FormEncoder {
	return form.NewEncoder()
}
