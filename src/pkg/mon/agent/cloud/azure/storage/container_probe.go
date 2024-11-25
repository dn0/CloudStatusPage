package storage

import (
	"context"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	containerProbeName         = "azure_storage_container" // doc="Azure Blob Storage Container"
	containerProbeActionCreate = 10                        // name="storage.container.create" doc="Creates a new storage container" url="https://learn.microsoft.com/en-us/rest/api/storageservices/create-container"
	containerProbeActionDelete = 20                        // name="storage.container.delete" doc="Deletes the container created by `storage.container.create`" url="https://learn.microsoft.com/en-us/rest/api/storageservices/delete-container"
)

type ContainerProbe[T agent.Azure] struct {
	container *container
}

func NewContainerProbe[T agent.Azure](cfg *agent.AzureConfig, storageClient *azblob.Client) *ContainerProbe[T] {
	return &ContainerProbe[T]{
		container: &container{
			namePrefix: cfg.Cloud.StorageContainerPrefix,
			client:     storageClient,
			timeout:    cfg.ProbeTimeout,
		},
	}
}

func (p *ContainerProbe[T]) String() string {
	return containerProbeName
}

func (p *ContainerProbe[T]) Start(_ context.Context) {
	p.log().Info(
		"Probe initialized",
		"storage_url", p.container.client.URL(),
		"container_prefix", p.container.namePrefix,
	)
}

func (p *ContainerProbe[T]) Do(ctx context.Context) []*pb.Result {
	var res *pb.Result

	if p.container.exists {
		res = pb.NewResult(containerProbeActionDelete)
		res.Store(p.deleteContainer(ctx))
		if res.Failed() {
			p.container.exists = false // let's start fresh
		}
	} else {
		p.container.new()
		res = pb.NewResult(containerProbeActionCreate)
		res.Store(p.createContainer(ctx))
	}

	return []*pb.Result{res}
}

func (p *ContainerProbe[T]) Stop(ctx context.Context) {
	if p.container.name != "" {
		p.log().Debug("Deleting storage container")
		_ = p.container.delete(ctx)
	}
}

func (p *ContainerProbe[T]) createContainer(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Creating storage container...")
	ret, err := pb.Timeit(p.container.create, ctx)
	if err == nil {
		p.log().Debug("Created storage container", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not create storage container", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *ContainerProbe[T]) deleteContainer(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Deleting storage container...", "probe", p.String(), "container", p.container.String())
	ret, err := pb.Timeit(p.container.delete, ctx)
	if err == nil {
		p.log().Debug("Deleted storage container", "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not delete storage container", "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *ContainerProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "container", p.container.String())
}
