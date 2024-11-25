package ec2

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"cspage/pkg/mon/agent"
	VPC "cspage/pkg/mon/agent/cloud/common/vpc"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	vmProbeName             = "aws_ec2_vm"      // doc="Amazon EC2 VM Instance"
	vmSpotProbeName         = "aws_ec2_vm_spot" // doc="Amazon EC2 VM Spot Instance"
	vmProbeActionCreate     = 10                // name="compute.vm.create"      doc="Creates a new vm instance" url="https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_RunInstances.html"
	vmProbeActionCreateWait = 11                // name="compute.vm.createWait"  doc="Waits for the new vm instance created by `compute.vm.create` to start" url=""
	vmProbeActionGet        = 20                // name="compute.vm.get"         doc="Describes the vm instance created by `compute.vm.create`" url="https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html"
	vmProbeActionDelete     = 30                // name="compute.vm.delete"      doc="Deletes the vm instance created by `compute.vm.create`" url="https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_TerminateInstances.html"
	vmProbeActionDeleteWait = 31                // name="compute.vm.deleteWait"  doc="Waits for the vm instance deleted by `compute.vm.delete` to get completely removed" url=""
	vmProbeStartDelayDiv    = 2
)

type VMProbe[T agent.AWS] struct {
	cfg        *agent.AWSConfig
	name       string
	vm         *vm
	pingJob    *VPC.PingProbeJob[agent.AWS]
	startDelay time.Duration
}

func NewVMProbe[T agent.AWS](
	cfg *agent.AWSConfig,
	awsConfig *aws.Config,
	model VMProvisioningModel,
	pingJob *VPC.PingProbeJob[agent.AWS],
) *VMProbe[T] {
	var name, vmPrefix, vmType, vmImage string
	var zonesSkip []string
	var startDelay time.Duration
	switch model {
	case VMProvisioningModelStandard:
		name = vmProbeName
		zonesSkip = cfg.Cloud.EC2VMZonesSkip
		vmPrefix = cfg.Cloud.EC2VMPrefix
		vmType = cfg.Cloud.EC2VMType
		vmImage = cfg.Cloud.EC2VMDiskImageName
		startDelay = cfg.ProbeLongIntervalDefault / vmProbeStartDelayDiv
	case VMProvisioningModelSpot:
		name = vmSpotProbeName
		zonesSkip = cfg.Cloud.EC2VMSpotZonesSkip
		vmPrefix = cfg.Cloud.EC2VMSpotPrefix
		vmType = cfg.Cloud.EC2VMSpotType
		vmImage = cfg.Cloud.EC2VMSpotDiskImageName
		startDelay = 0
	}

	return &VMProbe[T]{
		cfg:     cfg,
		name:    name,
		pingJob: pingJob,
		vm: &vm{
			namePrefix:        vmPrefix,
			zonesSkip:         zonesSkip,
			provisioningModel: model,
			machineType:       vmType,
			imageName:         vmImage,
			client:            ec2.NewFromConfig(*awsConfig),
			timeout:           cfg.ProbeTimeout,
			waitTimeout:       cfg.ProbeLongTimeout,
		},
		startDelay: startDelay,
	}
}

func (p *VMProbe[T]) String() string {
	return p.name
}

func (p *VMProbe[T]) Start(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, p.cfg.ProbeTimeout)
	defer cancel()
	var err error

	if p.vm.zones, err = getAvailabilityZones(ctx, p.vm.client, p.vm.zonesSkip); err != nil {
		agent.DieLog(p.log(), "Could not fetch AWS availability zones", "region", p.cfg.Env.Region, "err", err)
	}

	if p.vm.subnets, err = getVPCSubnets(
		ctx,
		p.vm.client,
		types.Filter{Name: aws.String("vpc-id"), Values: []string{p.cfg.Cloud.EC2VMVPCID}},
	); err != nil {
		agent.DieLog(p.log(), "Could not fetch AWS VPC subnets", "vpc_id", p.cfg.Cloud.EC2VMVPCID, "err", err)
	}

	for _, zone := range p.vm.zones {
		if _, ok := p.vm.subnets[zone]; !ok {
			agent.DieLog(p.log(), "Could not find VPC subnets for zone", "zone", zone, "vpc_id", p.cfg.Cloud.EC2VMVPCID)
		}
	}

	if p.vm.imageID, err = getLatestAMI(
		ctx,
		p.vm.client,
		types.Filter{Name: aws.String("owner-alias"), Values: []string{p.cfg.Cloud.EC2VMDiskImageOwner}},
		types.Filter{Name: aws.String("name"), Values: []string{p.vm.imageName}},
	); err != nil {
		agent.DieLog(p.log(), "Could not fetch latest AWS image ID", "region", p.cfg.Env.Region, "err", err)
	}

	p.log().Info("Probe start is delayed", "sleep_time", p.startDelay)
	time.Sleep(p.startDelay)

	p.log().Info(
		"Probe initialized",
		"zones", p.vm.zones,
		"zones_skip", p.vm.zonesSkip,
		"ami", p.vm.imageID,
		"vm_prefix", p.vm.namePrefix,
	)

	p.pingJob.Start(ctx)
}

func (p *VMProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(vmProbeActionCreate),
		pb.NewResult(vmProbeActionCreateWait),
		pb.NewResult(vmProbeActionGet),
		pb.NewResult(vmProbeActionDelete),
		pb.NewResult(vmProbeActionDeleteWait),
	}
	p.vm.new()

	p.createVM(ctx, res[0], res[1])
	if res[0].Failed() || res[1].Failed() {
		p.cleanup(ctx)
		return res
	}

	p.getVM(ctx, res[2])
	if res[2].Succeeded() {
		// TODO: whose problem is it?
		_ = p.pingJob.Run(ctx, p.vm.name, p.vm.id)
	}

	p.deleteVM(ctx, res[3], res[4])

	return res
}

func (p *VMProbe[T]) Stop(ctx context.Context) {
	p.pingJob.Stop(ctx)
	p.cleanup(ctx)
}

func (p *VMProbe[T]) cleanup(ctx context.Context) {
	if p.vm.id != "" {
		p.log().Debug("Deleting vm instance")
		_ = p.vm.delete(ctx)
	}
}

func (p *VMProbe[T]) createVM(ctx context.Context, resCreate, resWait *pb.Result) {
	p.log().Debug("Creating vm instance...")
	var err error
	if err = resCreate.Timeit(p.vm.create, ctx); err == nil {
		if err = resWait.Timeit(p.vm.createWait, ctx); err != nil {
			_ = p.vm.delete(ctx)
		}
	}
	if err == nil {
		p.log().Debug("Created vm instance", "took_create", resCreate.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not create vm instance", "took_create", resCreate.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *VMProbe[T]) getVM(ctx context.Context, res *pb.Result) {
	p.log().Debug("Getting vm instance...")
	err := res.Timeit(p.vm.get, ctx)
	if err == nil {
		p.log().Debug("Got vm instance", "took", res.Took())
	} else {
		p.log().Error("Could not get vm instance", "took", res.Took(), "err", err)
	}
}

func (p *VMProbe[T]) deleteVM(ctx context.Context, resDelete, resWait *pb.Result) {
	p.log().Debug("Deleting vm instance...")
	var err error
	if err = resDelete.Timeit(p.vm.delete, ctx); err == nil {
		err = resWait.Timeit(p.vm.deleteWait, ctx)
	}
	if err == nil {
		p.log().Debug("Deleted vm instance", "took_delete", resDelete.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not delete vm instance", "took_delete", resDelete.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *VMProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "vm", p.vm.String())
}
