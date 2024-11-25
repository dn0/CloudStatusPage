package compute

import (
	"context"
	"log/slog"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"

	"cspage/pkg/mon/agent"
	VPC "cspage/pkg/mon/agent/cloud/common/vpc"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	vmProbeName             = "azure_compute_vm"      // doc="Azure Virtual Machine"
	vmSpotProbeName         = "azure_compute_vm_spot" // doc="Azure Spot Virtual Machine"
	vmProbeActionCreate     = 10                      // name="compute.vm.create"      doc="Creates a new vm instance" url="https://learn.microsoft.com/en-us/rest/api/compute/virtual-machines/create-or-update"
	vmProbeActionCreateWait = 11                      // name="compute.vm.createWait"  doc="Waits for the new vm instance created by `compute.vm.create` to start" url=""
	vmProbeActionGet        = 20                      // name="compute.vm.get"         doc="Describes the vm instance created by `compute.vm.create`" url="https://learn.microsoft.com/en-us/rest/api/compute/virtual-machines/get"
	vmProbeActionDelete     = 30                      // name="compute.vm.delete"      doc="Deletes the vm instance created by `compute.vm.create`" url="https://learn.microsoft.com/en-us/rest/api/compute/virtual-machines/delete"
	vmProbeActionDeleteWait = 31                      // name="compute.vm.deleteWait"  doc="Waits for the vm instance deleted by `compute.vm.delete` to get completely removed" url=""
	vmProbeStartDelayDiv    = 2
)

type VMProbe[T agent.Azure] struct {
	cfg        *agent.AzureConfig
	name       string
	vm         *vm
	pingJob    *VPC.PingProbeJob[agent.Azure]
	startDelay time.Duration
}

func NewVMProbe[T agent.Azure](
	cfg *agent.AzureConfig,
	factory *armcompute.ClientFactory,
	model VMProvisioningModel,
	pingJob *VPC.PingProbeJob[agent.Azure],
) *VMProbe[T] {
	var name, vmPrefix, vmType, nicID string
	var startDelay time.Duration
	switch model {
	case VMProvisioningModelStandard:
		name = vmProbeName
		vmPrefix = cfg.Cloud.ComputeVMPrefix
		vmType = cfg.Cloud.ComputeVMType
		nicID = cfg.Cloud.ComputeVMNICID
		startDelay = cfg.ProbeLongIntervalDefault / vmProbeStartDelayDiv
	case VMProvisioningModelSpot:
		name = vmSpotProbeName
		vmPrefix = cfg.Cloud.ComputeVMSpotPrefix
		vmType = cfg.Cloud.ComputeVMSpotType
		nicID = cfg.Cloud.ComputeVMSpotNICID
		startDelay = 0
	}

	return &VMProbe[T]{
		cfg:     cfg,
		name:    name,
		pingJob: pingJob,
		vm: &vm{
			namePrefix:        vmPrefix,
			resourceGroup:     cfg.Cloud.ResourceGroup,
			provisioningModel: model,
			machineType:       vmType,
			nicID:             nicID,
			client:            factory.NewVirtualMachinesClient(),
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
	p.log().Info("Probe start is delayed", "sleep_time", p.startDelay)
	time.Sleep(p.startDelay)

	p.vm.init(
		p.cfg.Env.Region,
		p.cfg.Cloud.ComputeVMDiskImagePublisher,
		p.cfg.Cloud.ComputeVMDiskImageOffer,
		p.cfg.Cloud.ComputeVMDiskImageSKU,
		p.cfg.Cloud.ComputeVMDiskImageVersion,
	)
	p.log().Info(
		"Probe initialized",
		"resource_group", p.vm.resourceGroup,
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
