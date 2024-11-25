package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

const (
	blobProbeDownloadBufferSize = 1024
)

type blob struct {
	namePrefix    string
	containerName string
	client        *azblob.Client
	timeout       time.Duration

	name   string
	exists bool
}

func (b *blob) String() string {
	return "uri=" + b.client.URL() + "/" + b.containerName + "/" + b.name
}

func (b *blob) new() {
	now := time.Now()
	b.name = fmt.Sprintf("%s%d", b.namePrefix, now.UnixNano())
	b.exists = false
}

func (b *blob) upload(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	buf := []byte(b.String())
	if _, err := b.client.UploadBuffer(ctx, b.containerName, b.name, buf, nil); err != nil {
		return fmt.Errorf("blob(%s).Upload: %w", b, err)
	}

	b.exists = true
	return nil
}

func (b *blob) download(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	buf := make([]byte, blobProbeDownloadBufferSize)
	if _, err := b.client.DownloadBuffer(ctx, b.containerName, b.name, buf, nil); err != nil {
		return fmt.Errorf("blob(%s).Download: %w", b, err)
	}

	return nil
}

func (b *blob) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	if _, err := b.client.DeleteBlob(ctx, b.containerName, b.name, nil); err != nil {
		return fmt.Errorf("blob(%s).Delete: %w", b, err)
	}

	b.exists = false
	return nil
}
