package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	cloudStorage "cloud.google.com/go/storage"
)

const (
	objectProbeChunkSize = 0 // See docs for more info
)

type object struct {
	namePrefix string
	bucketName string
	bucket     *cloudStorage.BucketHandle
	timeout    time.Duration

	name   string
	exists bool
}

func (o *object) String() string {
	return "uri=gs://" + o.bucketName + "/" + o.name
}

func (o *object) new() {
	now := time.Now()
	o.name = fmt.Sprintf("%s%d", o.namePrefix, now.UnixNano())
	o.exists = false
}

func (o *object) upload(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	obj := o.bucket.Object(o.name)
	obj = obj.If(cloudStorage.Conditions{DoesNotExist: true})
	writer := obj.NewWriter(ctx)
	writer.ChunkSize = objectProbeChunkSize
	if _, err := fmt.Fprint(writer, o.String()); err != nil {
		_ = writer.Close()
		return fmt.Errorf("object(%s).Write: %w", o, err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("object(%s).WriterClose: %w", o, err)
	}

	o.exists = true
	return nil
}

func (o *object) download(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	obj := o.bucket.Object(o.name)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("object(%s).NewReader: %w", o, err)
	}
	if _, err := io.ReadAll(reader); err != nil {
		return fmt.Errorf("object(%s).ReadAll: %w", o, err)
	}

	return nil
}

func (o *object) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	obj := o.bucket.Object(o.name)

	// Set a generation-match precondition to avoid potential race conditions and data corruptions.
	// The request to delete the file is aborted if the object's generation number does not match your precondition.
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object(%s).Attrs: %w", o, err)
	}
	obj = obj.If(cloudStorage.Conditions{GenerationMatch: attrs.Generation})

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("object(%s).Delete: %w", o, err)
	}

	o.exists = false
	return nil
}
