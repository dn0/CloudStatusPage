package storage

import (
	"context"
	"log/slog"
	"time"

	cloudStorage "cloud.google.com/go/storage"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	objectProbeSleep          = 500 * time.Millisecond
	objectProbeName           = "gcp_storage_object" // doc="Google Cloud Storage Object"
	objectProbeActionUpload   = 10                   // name="storage.object.upload" doc="Uploads a small object to a bucket" url="https://cloud.google.com/storage/docs/json_api/v1/objects/insert"
	objectProbeActionDownload = 20                   // name="storage.object.download" doc="Downloads the object created by `storage.object.upload`" url="https://cloud.google.com/storage/docs/json_api/v1/objects/get"
	objectProbeActionDelete   = 30                   // name="storage.object.delete" doc="Deletes the object created by `storage.object.upload`" url="https://cloud.google.com/storage/docs/json_api/v1/objects/delete"
)

type ObjectProbe[T agent.GCP] struct {
	cfg    *agent.GCPConfig
	object *object
}

func NewObjectProbe[T agent.GCP](cfg *agent.GCPConfig, storageClient *cloudStorage.Client) *ObjectProbe[T] {
	return &ObjectProbe[T]{
		cfg: cfg,
		object: &object{
			namePrefix: cfg.Cloud.StorageObjectPrefix,
			bucketName: cfg.Cloud.StorageObjectBucketName,
			bucket:     storageClient.Bucket(cfg.Cloud.StorageObjectBucketName),
			timeout:    cfg.ProbeTimeout,
		},
	}
}

func (p *ObjectProbe[T]) String() string {
	return objectProbeName
}

func (p *ObjectProbe[T]) Start(_ context.Context) {
	p.log().Info(
		"Probe initialized",
		"bucket", p.object.bucket.BucketName(),
		"object_prefix", p.object.namePrefix,
	)
}

func (p *ObjectProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(objectProbeActionUpload),
		pb.NewResult(objectProbeActionDownload),
		pb.NewResult(objectProbeActionDelete),
	}
	p.object.new()

	res[0].Store(p.uploadObject(ctx))
	if res[0].Failed() {
		return res
	}
	time.Sleep(objectProbeSleep)
	res[1].Store(p.downloadObject(ctx))
	res[2].Store(p.deleteObject(ctx))

	return res
}

func (p *ObjectProbe[T]) Stop(ctx context.Context) {
	if p.object.name != "" {
		p.log().Debug("Deleting storage object")
		_ = p.object.delete(ctx)
	}
}

func (p *ObjectProbe[T]) uploadObject(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Uploading storage object...")
	ret, err := pb.Timeit(p.object.upload, ctx)
	if err == nil {
		p.log().Debug("Uploaded storage object", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not upload storage object", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *ObjectProbe[T]) downloadObject(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Downloading storage object...")
	ret, err := pb.Timeit(p.object.download, ctx)
	if err == nil {
		p.log().Debug("Downloaded storage object", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not download storage object", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *ObjectProbe[T]) deleteObject(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Deleting storage object...")
	ret, err := pb.Timeit(p.object.delete, ctx)
	if err == nil {
		p.log().Debug("Deleted storage object", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not delete storage object", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *ObjectProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "object", p.object.String())
}
