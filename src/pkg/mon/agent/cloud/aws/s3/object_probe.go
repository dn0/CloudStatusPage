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
	objectProbeName             = "aws_s3_object" // doc="Amazon S3 Object"
	objectProbeActionUpload     = 10              // name="s3.object.upload"      doc="Uploads a small object to a bucket" url="https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html"
	objectProbeActionUploadWait = 11              // name="s3.object.uploadWait"  doc="Waits for the object created by `s3.object.upload` to become available" url=""
	objectProbeActionDownload   = 20              // name="s3.object.download"    doc="Downloads the object created by `s3.object.upload`" url="https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html"
	objectProbeActionDelete     = 30              // name="s3.object.delete"      doc="Deletes the object created by `s3.object.upload`" url="https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html"
	objectProbeActionDeleteWait = 31              // name="s3.object.deleteWait"  doc="Waits for the object deleted by `s3.object.delete` to get completely removed" url=""
)

type ObjectProbe[T agent.AWS] struct {
	object *object
}

func NewObjectProbe[T agent.AWS](cfg *agent.AWSConfig, awsConfig *aws.Config) *ObjectProbe[T] {
	return &ObjectProbe[T]{
		object: &object{
			namePrefix:  cfg.Cloud.S3ObjectPrefix,
			bucketName:  cfg.Cloud.S3ObjectBucketName,
			client:      s3.NewFromConfig(*awsConfig),
			timeout:     cfg.ProbeTimeout,
			waitTimeout: cfg.ProbeTimeout,
		},
	}
}

func (p *ObjectProbe[T]) String() string {
	return objectProbeName
}

func (p *ObjectProbe[T]) Start(_ context.Context) {
	p.log().Info(
		"Probe initialized",
		"bucket", p.object.bucketName,
		"object_prefix", p.object.namePrefix,
	)
}

func (p *ObjectProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(objectProbeActionUpload),
		pb.NewResult(objectProbeActionUploadWait),
		pb.NewResult(objectProbeActionDownload),
		pb.NewResult(objectProbeActionDelete),
		pb.NewResult(objectProbeActionDeleteWait),
	}
	p.object.new()

	p.uploadObject(ctx, res[0], res[1])
	if res[0].Failed() || res[1].Failed() {
		return res
	}
	p.downloadObject(ctx, res[2])
	p.deleteObject(ctx, res[3], res[4])

	return res
}

func (p *ObjectProbe[T]) Stop(ctx context.Context) {
	if p.object.name != "" {
		p.log().Debug("Deleting S3 object")
		_ = p.object.delete(ctx)
	}
}

func (p *ObjectProbe[T]) uploadObject(ctx context.Context, resCreate, resWait *pb.Result) {
	p.log().Debug("Uploading S3 object...")
	var err error
	if err = resCreate.Timeit(p.object.upload, ctx); err == nil {
		if err = resWait.Timeit(p.object.uploadWait, ctx); err != nil {
			p.log().Debug("Deleting S3 object")
			_ = p.object.delete(ctx)
		}
	}
	if err == nil {
		p.log().Debug("Uploaded S3 object", "took_create", resCreate.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not upload S3 object", "took_create", resCreate.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *ObjectProbe[T]) downloadObject(ctx context.Context, res *pb.Result) {
	p.log().Debug("Downloading S3 object...")
	err := res.Timeit(p.object.download, ctx)
	if err == nil {
		p.log().Debug("Downloaded S3 object", "took", res.Took())
	} else {
		p.log().Error("Could not download S3 object", "took", res.Took(), "err", err)
	}
}

func (p *ObjectProbe[T]) deleteObject(ctx context.Context, resDelete, resWait *pb.Result) {
	p.log().Debug("Deleting S3 object...")
	var err error
	if err = resDelete.Timeit(p.object.delete, ctx); err == nil {
		err = resWait.Timeit(p.object.deleteWait, ctx)
	}
	if err == nil {
		p.log().Debug("Deleted S3 object", "took_delete", resDelete.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not delete S3 object", "took_delete", resDelete.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *ObjectProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "object", p.object.String())
}
