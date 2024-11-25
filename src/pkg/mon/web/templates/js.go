package templates

import (
	"context"
	"io"
)

type jsBytes []byte

func (js jsBytes) Render(_ context.Context, w io.Writer) error {
	_, err := w.Write(js)
	//nolint:wrapcheck // There isn't much to wrap.
	return err
}
