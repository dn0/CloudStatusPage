package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type container struct {
	namePrefix string
	client     *azblob.Client
	timeout    time.Duration

	name   string
	exists bool
}

func (c *container) String() string {
	return "uri=" + c.client.URL() + "/" + c.name
}

func (c *container) new() {
	c.name = fmt.Sprintf("%s-%d", c.namePrefix, time.Now().Unix())
}

func (c *container) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if _, err := c.client.CreateContainer(ctx, c.name, nil); err != nil {
		return fmt.Errorf("container(%s).Create: %w", c, err)
	}

	c.exists = true
	return nil
}

func (c *container) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if _, err := c.client.DeleteContainer(ctx, c.name, nil); err != nil {
		return fmt.Errorf("container(%s).Delete: %w", c, err)
	}

	c.exists = false
	return nil
}
