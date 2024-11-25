package storage

import (
	"context"
	"fmt"
	"time"

	cloudStorage "cloud.google.com/go/storage"
)

type bucket struct {
	namePrefix string
	projectID  string
	attrs      *cloudStorage.BucketAttrs
	client     *cloudStorage.Client
	timeout    time.Duration

	name   string
	exists bool
}

func (b *bucket) String() string {
	return "uri=gs://" + b.name
}

func (b *bucket) new() {
	b.name = fmt.Sprintf("%s-%d", b.namePrefix, time.Now().Unix())
}

func (b *bucket) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	bucket := b.client.Bucket(b.name)
	if err := bucket.Create(ctx, b.projectID, b.attrs); err != nil {
		return fmt.Errorf("bucket(%s).Create: %w", b, err)
	}

	b.exists = true
	return nil
}

func (b *bucket) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	bucket := b.client.Bucket(b.name)
	if err := bucket.Delete(ctx); err != nil {
		return fmt.Errorf("bucket(%s).Delete: %w", b, err)
	}

	b.exists = false
	return nil
}
