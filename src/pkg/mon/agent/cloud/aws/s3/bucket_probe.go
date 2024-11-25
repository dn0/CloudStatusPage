package s3

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	bucketProbeName             = "aws_s3_bucket" // doc="Amazon S3 Bucket"
	bucketProbeActionCreate     = 10              // name="s3.bucket.create"      doc="Creates a new S3 bucket" url="https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html"
	bucketProbeActionCreateWait = 11              // name="s3.bucket.createWait"  doc="Waits for the bucket created by `s3.bucket.create` to become available" url=""
	bucketProbeActionDelete     = 20              // name="s3.bucket.delete"      doc="Deletes the bucket created by `s3.bucket.create`" url="https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html"
	bucketProbeActionDeleteWait = 21              // name="s3.bucket.deleteWait"  doc="Waits for the bucket delete by `s3.bucket.delete` to get completely removed" url=""
)

type BucketProbe[T agent.AWS] struct {
	bucket *bucket
}

func NewBucketProbe[T agent.AWS](cfg *agent.AWSConfig, awsConfig *aws.Config) *BucketProbe[T] {
	return &BucketProbe[T]{
		bucket: &bucket{
			namePrefix:  cfg.Cloud.S3BucketPrefix,
			region:      cfg.Env.Region,
			client:      s3.NewFromConfig(*awsConfig),
			timeout:     cfg.ProbeTimeout,
			waitTimeout: cfg.ProbeTimeout,
		},
	}
}

func (p *BucketProbe[T]) String() string {
	return bucketProbeName
}

func (p *BucketProbe[T]) Start(_ context.Context) {
	p.log().Info(
		"Probe initialized",
		"bucket_prefix", p.bucket.namePrefix,
	)
}

func (p *BucketProbe[T]) Do(ctx context.Context) []*pb.Result {
	var res1, res2 *pb.Result

	if p.bucket.exists {
		res1 = pb.NewResult(bucketProbeActionDelete)
		res2 = pb.NewResult(bucketProbeActionDeleteWait)
		p.deleteBucket(ctx, res1, res2)
		if res1.Failed() || res2.Failed() {
			p.bucket.exists = false // let's start fresh
		}
	} else {
		p.bucket.new()
		res1 = pb.NewResult(bucketProbeActionCreate)
		res2 = pb.NewResult(bucketProbeActionCreateWait)
		p.createBucket(ctx, res1, res2)
	}

	return []*pb.Result{res1, res2}
}

func (p *BucketProbe[T]) Stop(ctx context.Context) {
	if p.bucket.name != "" {
		p.log().Debug("Deleting S3 bucket")
		_ = p.bucket.delete(ctx)
	}
}

func (p *BucketProbe[T]) createBucket(ctx context.Context, resCreate, resWait *pb.Result) {
	p.log().Debug("Creating S3 bucket...")
	var err error
	if err = resCreate.Timeit(p.bucket.create, ctx); err == nil {
		if err = resWait.Timeit(p.bucket.createWait, ctx); err != nil {
			p.log().Debug("Deleting S3 bucket")
			_ = p.bucket.delete(ctx)
		}
	}
	if err == nil {
		p.log().Debug("Created S3 bucket", "took_create", resCreate.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not create S3 bucket", "took_create", resCreate.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *BucketProbe[T]) deleteBucket(ctx context.Context, resDelete, resWait *pb.Result) {
	p.log().Debug("Deleting S3 bucket...")
	var err error
	if err = resDelete.Timeit(p.bucket.delete, ctx); err == nil {
		err = resWait.Timeit(p.bucket.deleteWait, ctx)
	}
	if err == nil {
		p.log().Debug("Deleted S3 bucket", "took_delete", resDelete.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not delete S3 bucket", "took_delete", resDelete.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *BucketProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "bucket", p.bucket.String())
}
