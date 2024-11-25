package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	blobProbeSleep          = 500 * time.Millisecond
	blobProbeName           = "azure_storage_blob" // doc="Azure Blob Storage Object"
	blobProbeActionUpload   = 10                   // name="storage.blob.upload" doc="Uploads a small blob to a storage container" url="https://learn.microsoft.com/en-us/rest/api/storageservices/put-blob"
	blobProbeActionDownload = 20                   // name="storage.blob.download" doc="Downloads the blob created by `storage.blob.upload`" url="https://learn.microsoft.com/en-us/rest/api/storageservices/get-blob"
	blobProbeActionDelete   = 30                   // name="storage.blob.delete" doc="Deletes the blob created by `storage.blob.upload`" url="https://learn.microsoft.com/en-us/rest/api/storageservices/delete-blob"
)

type BlobProbe[T agent.Azure] struct {
	blob *blob
}

func NewBlobProbe[T agent.Azure](cfg *agent.AzureConfig, storageClient *azblob.Client) *BlobProbe[T] {
	return &BlobProbe[T]{
		blob: &blob{
			namePrefix:    cfg.Cloud.StorageBlobPrefix,
			containerName: cfg.Cloud.StorageBlobContainerName,
			client:        storageClient,
			timeout:       cfg.ProbeTimeout,
		},
	}
}

func (p *BlobProbe[T]) String() string {
	return blobProbeName
}

func (p *BlobProbe[T]) Start(_ context.Context) {
	p.log().Info(
		"Probe initialized",
		"storage_url", p.blob.client.URL(),
		"container", p.blob.containerName,
		"blob_prefix", p.blob.namePrefix,
	)
}

func (p *BlobProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(blobProbeActionUpload),
		pb.NewResult(blobProbeActionDownload),
		pb.NewResult(blobProbeActionDelete),
	}
	p.blob.new()

	res[0].Store(p.uploadBlob(ctx))
	if res[0].Failed() {
		return res
	}
	time.Sleep(blobProbeSleep)
	res[1].Store(p.downloadBlob(ctx))
	res[2].Store(p.deleteBlob(ctx))

	return res
}

func (p *BlobProbe[T]) Stop(ctx context.Context) {
	if p.blob.name != "" {
		p.log().Debug("Deleting storage blob")
		_ = p.blob.delete(ctx)
	}
}

func (p *BlobProbe[T]) uploadBlob(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Uploading storage blob...")
	ret, err := pb.Timeit(p.blob.upload, ctx)
	if err == nil {
		p.log().Debug("Uploaded storage blob", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not upload storage blob", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *BlobProbe[T]) downloadBlob(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Downloading storage blob...")
	ret, err := pb.Timeit(p.blob.download, ctx)
	if err == nil {
		p.log().Debug("Downloaded storage blob", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not download storage blob", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *BlobProbe[T]) deleteBlob(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Deleting storage blob...")
	ret, err := pb.Timeit(p.blob.delete, ctx)
	if err == nil {
		p.log().Debug("Deleted storage blob", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not delete storage blob", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *BlobProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "blob", p.blob.String())
}
