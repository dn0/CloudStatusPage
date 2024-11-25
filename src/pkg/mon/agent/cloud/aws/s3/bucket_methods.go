package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	regionUSEast1 = "us-east-1"
)

type bucket struct {
	namePrefix  string
	region      string
	client      *s3.Client
	timeout     time.Duration
	waitTimeout time.Duration

	name   string
	exists bool

	createBucketConfiguration *types.CreateBucketConfiguration
}

func (b *bucket) String() string {
	return "uri=s3://" + b.name
}

func (b *bucket) new() {
	b.name = fmt.Sprintf("%s-%d", b.namePrefix, time.Now().Unix())
	// AWS Docs: If you are creating a bucket on the US East (N. Virginia) region (us-east-1),
	//           you do not need to specify the location constraint.
	if b.region == regionUSEast1 {
		b.createBucketConfiguration = nil
	} else {
		b.createBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(b.region),
		}
	}
}

func (b *bucket) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	if _, err := b.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket:                    aws.String(b.name),
		CreateBucketConfiguration: b.createBucketConfiguration,
	}); err != nil {
		return fmt.Errorf("bucket(%s).CreateBucket: %w", b, err)
	}

	b.exists = true
	return nil
}

func (b *bucket) createWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	waiter := s3.NewBucketExistsWaiter(b.client, setBucketExistsWaiterOptions)
	if err := waiter.Wait(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(b.name),
	}, b.waitTimeout); err != nil {
		return fmt.Errorf("bucket(%s).CreateBucket.wait: %w", b, err)
	}

	return nil
}

func (b *bucket) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	if _, err := b.client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(b.name),
	}); err != nil {
		return fmt.Errorf("bucket(%s).DeleteBucket: %w", b, err)
	}

	b.exists = false
	return nil
}

func (b *bucket) deleteWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	waiter := s3.NewBucketNotExistsWaiter(b.client, setBucketNotExistsWaiterOptions)
	if err := waiter.Wait(ctx, &s3.HeadBucketInput{Bucket: aws.String(b.name)}, b.waitTimeout); err != nil {
		return fmt.Errorf("bucket(%s).DeleteBucket.wait: %w", b, err)
	}

	return nil
}

func setBucketExistsWaiterOptions(options *s3.BucketExistsWaiterOptions) {
	options.MinDelay = waiterMinDelay
	options.MaxDelay = waiterMaxDelay
}

func setBucketNotExistsWaiterOptions(options *s3.BucketNotExistsWaiterOptions) {
	options.MinDelay = waiterMinDelay
	options.MaxDelay = waiterMaxDelay
}
