package s3

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type object struct {
	namePrefix  string
	bucketName  string
	client      *s3.Client
	timeout     time.Duration
	waitTimeout time.Duration

	name   string
	exists bool
}

func (o *object) String() string {
	return "uri=s3://" + o.bucketName + "/" + o.name
}

func (o *object) new() {
	now := time.Now()
	o.name = fmt.Sprintf("%s%d", o.namePrefix, now.UnixNano())
	o.exists = false
}

func (o *object) upload(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	if _, err := o.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(o.bucketName),
		Key:    aws.String(o.name),
		Body:   strings.NewReader(o.String()),
	}); err != nil {
		return fmt.Errorf("object(%s).PutObject: %w", o, err)
	}

	o.exists = true
	return nil
}

func (o *object) uploadWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	waiter := s3.NewObjectExistsWaiter(o.client, setObjectExistsWaiterOptions)
	if err := waiter.Wait(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(o.bucketName),
		Key:    aws.String(o.name),
	}, o.waitTimeout); err != nil {
		return fmt.Errorf("object(%s).PutObject.wait: %w", o, err)
	}

	return nil
}

func (o *object) download(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	res, err := o.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(o.bucketName),
		Key:    aws.String(o.name),
	})
	if err != nil {
		return fmt.Errorf("object(%s).GetObject: %w", o, err)
	}
	if _, err := io.ReadAll(res.Body); err != nil {
		_ = res.Body.Close()
		return fmt.Errorf("object(%s).ReadAll: %w", o, err)
	}
	if err := res.Body.Close(); err != nil {
		return fmt.Errorf("object(%s).BodyClose: %w", o, err)
	}

	return nil
}

func (o *object) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	if _, err := o.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(o.bucketName),
		Key:    aws.String(o.name),
	}); err != nil {
		return fmt.Errorf("object(%s).DeleteObjects: %w", o, err)
	}

	o.exists = false
	return nil
}

func (o *object) deleteWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	waiter := s3.NewObjectNotExistsWaiter(o.client, setObjectNotExistsWaiterOptions)
	if err := waiter.Wait(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(o.bucketName),
		Key:    aws.String(o.name),
	}, o.waitTimeout); err != nil {
		return fmt.Errorf("object(%s).DeleteObjects.wait: %w", o, err)
	}

	return nil
}

func setObjectExistsWaiterOptions(options *s3.ObjectExistsWaiterOptions) {
	options.MinDelay = waiterMinDelay
	options.MaxDelay = waiterMaxDelay
}

func setObjectNotExistsWaiterOptions(options *s3.ObjectNotExistsWaiterOptions) {
	options.MinDelay = waiterMinDelay
	options.MaxDelay = waiterMaxDelay
}
