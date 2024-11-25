package ec2

import (
	"context"
	"log/slog"

	"cspage/pkg/http"
	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	vmMetadataProbeName   = "aws_ec2_vm_metadata" // doc="Amazon EC2 VM Instance Metadata"
	vmMetadataProbeGet    = 10                    // name="compute.vm.metadata.get" doc="Get EC2 instance metadata" url="https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html"
	vmMetadataProbeSuffix = "tags/instance/app_version"
)

type VMMetadataProbe[T agent.AWS] struct {
	cfg    *agent.AWSConfig
	client *http.Client
}

func NewVMMetadataProbe[T agent.AWS](cfg *agent.AWSConfig) *VMMetadataProbe[T] {
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
