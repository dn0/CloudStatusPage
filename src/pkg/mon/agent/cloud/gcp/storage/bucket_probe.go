package storage

import (
	"context"
	"log/slog"

	cloudStorage "cloud.google.com/go/storage"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	bucketProbeName         = "gcp_storage_bucket" // doc="Google Cloud Storage Bucket"
	bucketProbeActionCreate = 10                   // name="storage.bucket.create" doc="Creates a new storage bucket" url="https://cloud.google.com/storage/docs/json_api/v1/buckets/insert"
	bucketProbeActionDelete = 20                   // name="storage.bucket.delete" doc="Deletes the bucket created by `storage.bucket.create`" url="https://cloud.google.com/storage/docs/json_api/v1/buckets/delete"
)

type BucketProbe[T agent.GCP] struct {
	cfg    *agent.GCPConfig
	bucket *bucket
}

func NewBucketProbe[T agent.GCP](cfg *agent.GCPConfig, storageClient *cloudStorage.Client) *BucketProbe[T] {
	return &BucketProbe[T]{
		cfg: cfg,
		bucket: &bucket{
			namePrefix: cfg.Cloud.StorageBucketPrefix,
			projectID:  cfg.Cloud.ProjectID,
			client:     storageClient,
			timeout:    cfg.ProbeTimeout,
			attrs: &cloudStorage.BucketAttrs{
				Location: cfg.Env.Region,
				Labels: map[string]string{
					"cost-center": "mon-probe",
				},
			},
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
	var res *pb.Result

	if p.bucket.exists {
		res = pb.NewResult(bucketProbeActionDelete)
		res.Store(p.deleteBucket(ctx))
		if res.Failed() {
			p.bucket.exists = false // let's start fresh
		}
	} else {
		p.bucket.new()
		res = pb.NewResult(bucketProbeActionCreate)
		res.Store(p.createBucket(ctx))
	}

	return []*pb.Result{res}
}

func (p *BucketProbe[T]) Stop(ctx context.Context) {
	if p.bucket.name != "" {
		p.log().Debug("Deleting storage bucket")
		_ = p.bucket.delete(ctx)
	}
}

func (p *BucketProbe[T]) createBucket(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Creating storage bucket...")
	ret, err := pb.Timeit(p.bucket.create, ctx)
	if err == nil {
		p.log().Debug("Created storage bucket", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not create storage bucket", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *BucketProbe[T]) deleteBucket(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Deleting storage bucket...")
	ret, err := pb.Timeit(p.bucket.delete, ctx)
	if err == nil {
		p.log().Debug("Deleted storage bucket", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not delete storage bucket", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *BucketProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "bucket", p.bucket.String())
}
