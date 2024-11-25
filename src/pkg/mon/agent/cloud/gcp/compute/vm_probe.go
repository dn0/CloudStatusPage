package compute

import (
	"context"
	"log/slog"
	"time"

	cloudCompute "cloud.google.com/go/compute/apiv1"

	"cspage/pkg/mon/agent"
	VPC "cspage/pkg/mon/agent/cloud/common/vpc"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	vmProbeName             = "gcp_compute_vm"      // doc="Google Cloud Virtual Machine"
	vmSpotProbeName         = "gcp_compute_vm_spot" // doc="Google Cloud Spot Virtual Machine"
	vmProbeActionCreate     = 10                    // name="compute.vm.create"      doc="Creates a new vm instance" url="https://cloud.google.com/compute/docs/reference/rest/v1/instances/insert"
	vmProbeActionCreateWait = 11                    // name="compute.vm.createWait"  doc="Waits for the new vm instance created by `compute.vm.create` to start" url="https://cloud.google.com/compute/docs/reference/rest/v1/zoneOperations/wait"
	vmProbeActionGet        = 20                    // name="compute.vm.get"         doc="Describes the vm instance created by `compute.vm.create`" url="https://cloud.google.com/compute/docs/reference/rest/v1/instances/get"
	vmProbeActionDelete     = 30                    // name="compute.vm.delete"      doc="Deletes the vm instance created by `compute.vm.create`" url="https://cloud.google.com/compute/docs/reference/rest/v1/instances/delete"
	vmProbeActionDeleteWait = 31                    // name="compute.vm.deleteWait"  doc="Waits for the vm instance deleted by `compute.vm.delete` to get completely removed" url="https://cloud.google.com/compute/docs/reference/rest/v1/zoneOperations/wait"
	vmProbeStartDelayDiv    = 2
)

type VMProbe[T agent.GCP] struct {
	cfg        *agent.GCPConfig
	name       string
	vm         *vm
	pingJob    *VPC.PingProbeJob[agent.GCP]
	startDelay time.Duration
}

func NewVMProbe[T agent.GCP](
	cfg *agent.GCPConfig,
	model VMProvisioningModel,
	pingJob *VPC.PingProbeJob[agent.GCP],
) *VMProbe[T] {
	var name, vmPrefix, vmType string
	var zonesSkip []string
	var startDelay time.Duration
	switch model {
	case VMProvisioningModelStandard:
		name = vmProbeName
		zonesSkip = cfg.Cloud.ComputeVMZonesSkip
		vmPrefix = cfg.Cloud.ComputeVMPrefix
		vmType = cfg.Cloud.ComputeVMType
		startDelay = cfg.ProbeLongIntervalDefault / vmProbeStartDelayDiv
	case VMProvisioningModelSpot:
		name = vmSpotProbeName
		zonesSkip = cfg.Cloud.ComputeVMSpotZonesSkip
		vmPrefix = cfg.Cloud.ComputeVMSpotPrefix
		vmType = cfg.Cloud.ComputeVMSpotType
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

	if p.vm.zones, err = getAvailabilityZones(ctx, p.cfg.Cloud.ProjectID, p.cfg.Env.Region, p.vm.zonesSkip); err != nil {
		agent.DieLog(p.log(), "Could not fetch GCP availability zones", "region", p.cfg.Env.Region, "err", err)
	}

	if p.vm.client, err = cloudCompute.NewInstancesRESTClient(ctx); err != nil {
		agent.DieLog(p.log(), "Could not initialize GCP vm instances client", "err", err)
	}

	p.log().Info("Probe start is delayed", "sleep_time", p.startDelay)
	time.Sleep(p.startDelay)

	p.vm.init(
		p.cfg.Cloud.ProjectID,
		p.cfg.Cloud.ComputeVMDiskImage,
		p.cfg.Cloud.ComputeVMSubnetwork,
	)
	p.log().Info(
		"Probe initialized",
		"zones", p.vm.zones,
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
		_ = p.pingJob.Run(ctx, p.vm.name, p.vm.name)
	}

	p.deleteVM(ctx, res[3], res[4])

	return res
}

func (p *VMProbe[T]) Stop(ctx context.Context) {
	p.pingJob.Stop(ctx)
	p.cleanup(ctx)

	if p.vm.client != nil {
		p.log().Debug("Closing GCP vm instances client...")
		if err := p.vm.client.Close(); err != nil {
			p.log().Error("Could not close GCP vm instances client", "err", err)
		}
	}
}

func (p *VMProbe[T]) cleanup(ctx context.Context) {
	if p.vm.name != "" {
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
