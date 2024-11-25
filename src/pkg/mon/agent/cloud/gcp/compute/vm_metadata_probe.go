package compute

import (
	"context"
	"log/slog"

	"cspage/pkg/http"
	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	vmMetadataProbeName   = "gcp_compute_vm_metadata" // doc="Google Cloud Virtual Machine Metadata"
	vmMetadataProbeGet    = 10                        // name="compute.vm.metadata.get" doc="Get virtual machine instance metadata" url="https://cloud.google.com/compute/docs/metadata"
	vmMetadataProbeSuffix = "/attributes/app_version"
)

type VMMetadataProbe[T agent.GCP] struct {
	cfg    *agent.GCPConfig
	client *http.Client
}

func NewVMMetadataProbe[T agent.GCP](cfg *agent.GCPConfig) *VMMetadataProbe[T] {
	return &VMMetadataProbe[T]{
		cfg:    cfg,
		client: http.NewClient(),
	}
}

func (p *VMMetadataProbe[T]) String() string {
	return vmMetadataProbeName
}

func (p *VMMetadataProbe[T]) Start(_ context.Context) {
	p.log().Info("Probe initialized")
}

func (p *VMMetadataProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := pb.NewResult(vmMetadataProbeGet)
	res.Store(p.getInstanceAttributeResult(ctx))

	return []*pb.Result{res}
}

func (p *VMMetadataProbe[T]) Stop(_ context.Context) {}

func (p *VMMetadataProbe[T]) getInstanceAttributeResult(ctx context.Context) (pb.ResultTime, error) {
	url := vmMetadataBaseURL + vmMetadataProbeSuffix
	p.log().Debug("Fetching vm metadata...", "url", url)
	ret, err := pb.Timeit(p.getInstanceAttribute, ctx)
	if err == nil {
		p.log().Debug("Got vm metadata", "url", url, "took", ret.Took)
	} else {
		//goland:noinspection GoDfaErrorMayBeNotNil
		p.log().Error("Could not fetch vm metadata", "url", url, "took", ret.Took, "err", err)
	}
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *VMMetadataProbe[T]) getInstanceAttribute(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, p.cfg.ProbeTimeout)
	defer cancel()

	if _, err := getInstanceMetadata(ctx, p.client, vmMetadataProbeSuffix); err != nil {
		return err
	}
	return nil
}

func (p *VMMetadataProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String())
}