package compute

import (
	"context"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	vhdSnapshotProbeName             = "azure_compute_vhd_snapshot" // doc="Azure Virtual Hard Disk Snapshot"
	vhdSnapshotProbeActionCreate     = 10                           // name="compute.vhd.snapshot.create"      doc="Creates a new VHD snapshot" url="https://learn.microsoft.com/en-us/rest/api/compute/snapshots/create-or-update"
	vhdSnapshotProbeActionCreateWait = 11                           // name="compute.vhd.snapshot.createWait"  doc="Waits for a new VHD snapshot created by `compute.vhd.snapshot.create` to complete" url=""
	vhdSnapshotProbeActionGet        = 20                           // name="compute.vhd.snapshot.get"         doc="Describes the VHD snapshot created by `compute.vhd.snapshot.create`" url="https://learn.microsoft.com/en-us/rest/api/compute/snapshots/get"
	vhdSnapshotProbeActionDelete     = 30                           // name="compute.vhd.snapshot.delete"      doc="Deletes the VHD snapshot created by `compute.vhd.snapshot.create`" url="https://learn.microsoft.com/en-us/rest/api/compute/snapshots/delete"
	vhdSnapshotProbeActionDeleteWait = 31                           // name="compute.vhd.snapshot.deleteWait"  doc="Waits for the VHD snapshot deleted by `compute.vhd.snapshot.delete` to get completely removed" url=""
)

type VHDSnapshotProbe[T agent.Azure] struct {
	cfg      *agent.AzureConfig
	snapshot *vhdSnapshot
}

func NewVHDSnapshotProbe[T agent.Azure](
	cfg *agent.AzureConfig,
	factory *armcompute.ClientFactory,
) *VHDSnapshotProbe[T] {
	return &VHDSnapshotProbe[T]{
		cfg: cfg,
		snapshot: &vhdSnapshot{
			namePrefix:    cfg.Cloud.ComputeVHDSnapshotPrefix,
			resourceGroup: cfg.Cloud.ResourceGroup,
			client:        factory.NewSnapshotsClient(),
			timeout:       cfg.ProbeTimeout,
			waitTimeout:   cfg.ProbeLongTimeout,
		},
	}
}

func (p *VHDSnapshotProbe[T]) String() string {
	return vhdSnapshotProbeName
}

func (p *VHDSnapshotProbe[T]) Start(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, p.cfg.ProbeTimeout)
	defer cancel()

	diskID := p.cfg.Cloud.ComputeVHDSnapshotDiskID
	if diskID == "" {
		var err error
		if diskID, err = getManagedDiskId(ctx); err != nil {
			agent.DieLog(p.log(), "Could not fetch managed disk ID", "err", err)
		}
	}

	p.snapshot.init(p.cfg.Env.Region, diskID)
	p.log().Info(
		"Probe initialized",
		"disk_id", diskID,
		"resource_group", p.snapshot.resourceGroup,
		"snapshot_prefix", p.snapshot.namePrefix,
	)
}

func (p *VHDSnapshotProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(vhdSnapshotProbeActionCreate),
		pb.NewResult(vhdSnapshotProbeActionCreateWait),
		pb.NewResult(vhdSnapshotProbeActionGet),
		pb.NewResult(vhdSnapshotProbeActionDelete),
		pb.NewResult(vhdSnapshotProbeActionDeleteWait),
	}
	p.snapshot.new()

	p.createVHDSnapshot(ctx, res[0], res[1])
	if res[0].Failed() || res[1].Failed() {
		p.cleanup(ctx)
		return res
	}
	p.getVHDSnapshot(ctx, res[2])
	p.deleteVHDSnapshot(ctx, res[3], res[4])

	return res
}

func (p *VHDSnapshotProbe[T]) Stop(ctx context.Context) {
	p.cleanup(ctx)
}

func (p *VHDSnapshotProbe[T]) cleanup(ctx context.Context) {
	if p.snapshot.name != "" {
		p.log().Debug("Deleting VHD snapshot")
		_ = p.snapshot.delete(ctx)
	}
}

func (p *VHDSnapshotProbe[T]) createVHDSnapshot(ctx context.Context, resCreate, resWait *pb.Result) {
	p.log().Debug("Creating VHD snapshot...")
	var err error
	if err = resCreate.Timeit(p.snapshot.create, ctx); err == nil {
		if err = resWait.Timeit(p.snapshot.createWait, ctx); err != nil {
			_ = p.snapshot.delete(ctx)
		}
	}
	if err == nil {
		p.log().Debug("Created VHD snapshot", "took_create", resCreate.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not create VHD snapshot", "took_create", resCreate.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *VHDSnapshotProbe[T]) getVHDSnapshot(ctx context.Context, res *pb.Result) {
	p.log().Debug("Getting VHD snapshot...")
	err := res.Timeit(p.snapshot.get, ctx)
	if err == nil {
		p.log().Debug("Got VHD snapshot", "took", res.Took())
	} else {
		p.log().Error("Could not get VHD snapshot", "took", res.Took(), "err", err)
	}
}

func (p *VHDSnapshotProbe[T]) deleteVHDSnapshot(ctx context.Context, resDelete, resWait *pb.Result) {
	p.log().Debug("Deleting VHD snapshot...")
	var err error
	if err = resDelete.Timeit(p.snapshot.delete, ctx); err == nil {
		err = resWait.Timeit(p.snapshot.deleteWait, ctx)
	}
	if err == nil {
		p.log().Debug("Deleted VHD snapshot", "took_delete", resDelete.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not delete VHD snapshot", "took_delete", resDelete.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *VHDSnapshotProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "snapshot", p.snapshot.String())
}
