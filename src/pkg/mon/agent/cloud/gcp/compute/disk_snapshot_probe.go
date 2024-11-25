package compute

import (
	"context"
	"log/slog"

	cloudCompute "cloud.google.com/go/compute/apiv1"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	diskSnapshotProbeName             = "gcp_compute_disk_snapshot" // doc="Google Cloud Disk Snapshot"
	diskSnapshotProbeActionCreate     = 10                          // name="compute.disk.snapshot.create"      doc="Creates a new disk snapshot" url="https://cloud.google.com/compute/docs/reference/rest/v1/snapshots/insert"
	diskSnapshotProbeActionCreateWait = 11                          // name="compute.disk.snapshot.createWait"  doc="Waits for a new disk snapshot created by `compute.disk.snapshot.create` to complete" url="https://cloud.google.com/compute/docs/reference/rest/v1/regionOperations/wait"
	diskSnapshotProbeActionGet        = 20                          // name="compute.disk.snapshot.get"         doc="Describes the disk snapshot created by `compute.disk.snapshot.create`" url="https://cloud.google.com/compute/docs/reference/rest/v1/snapshots/get"
	diskSnapshotProbeActionDelete     = 30                          // name="compute.disk.snapshot.delete"      doc="Deletes the disk snapshot created by `compute.disk.snapshot.create`" url="https://cloud.google.com/compute/docs/reference/rest/v1/snapshots/delete"
	diskSnapshotProbeActionDeleteWait = 31                          // name="compute.disk.snapshot.deleteWait"  doc="Waits for the disk snapshot deleted by `compute.disk.snapshot.delete` to get completely removed" url="https://cloud.google.com/compute/docs/reference/rest/v1/regionOperations/wait"
)

type DiskSnapshotProbe[T agent.GCP] struct {
	cfg      *agent.GCPConfig
	snapshot *diskSnapshot
}

func NewDiskSnapshotProbe[T agent.GCP](cfg *agent.GCPConfig) *DiskSnapshotProbe[T] {
	return &DiskSnapshotProbe[T]{
		cfg: cfg,
		snapshot: &diskSnapshot{
			namePrefix:  cfg.Cloud.ComputeDiskSnapshotPrefix,
			timeout:     cfg.ProbeTimeout,
			waitTimeout: cfg.ProbeLongTimeout,
		},
	}
}

func (p *DiskSnapshotProbe[T]) String() string {
	return diskSnapshotProbeName
}

func (p *DiskSnapshotProbe[T]) Start(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, p.cfg.ProbeTimeout)
	defer cancel()

	var err error
	if p.snapshot.client, err = cloudCompute.NewSnapshotsRESTClient(ctx); err != nil {
		agent.DieLog(p.log(), "Could not initialize GCP disk snapshot client", "err", err)
	}

	var diskName string
	if p.cfg.Cloud.ComputeDiskSnapshotDiskName == "" {
		diskName = p.cfg.Env.VM.Name
	} else {
		diskName = p.cfg.Cloud.ComputeDiskSnapshotDiskName
	}

	p.snapshot.init(p.cfg.Cloud.ProjectID, p.cfg.Env.Region, p.cfg.Env.Zone, diskName)
	p.log().Info(
		"Probe initialized",
		"disk_name", diskName,
		"snapshot_prefix", p.snapshot.namePrefix,
	)
}

func (p *DiskSnapshotProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(diskSnapshotProbeActionCreate),
		pb.NewResult(diskSnapshotProbeActionCreateWait),
		pb.NewResult(diskSnapshotProbeActionGet),
		pb.NewResult(diskSnapshotProbeActionDelete),
		pb.NewResult(diskSnapshotProbeActionDeleteWait),
	}
	p.snapshot.new()

	p.createDiskSnapshot(ctx, res[0], res[1])
	if res[0].Failed() || res[1].Failed() {
		p.cleanup(ctx)
		return res
	}
	p.getDiskSnapshot(ctx, res[2])
	p.deleteDiskSnapshot(ctx, res[3], res[4])

	return res
}

func (p *DiskSnapshotProbe[T]) Stop(ctx context.Context) {
	p.cleanup(ctx)

	if p.snapshot.client != nil {
		p.log().Debug("Closing GCP disk snapshot client...")
		if err := p.snapshot.client.Close(); err != nil {
			p.log().Error("Could not close GCP disk snapshot client", "err", err)
		}
	}
}

func (p *DiskSnapshotProbe[T]) cleanup(ctx context.Context) {
	if p.snapshot.name != "" {
		p.log().Debug("Deleting disk snapshot")
		_ = p.snapshot.delete(ctx)
	}
}

func (p *DiskSnapshotProbe[T]) createDiskSnapshot(ctx context.Context, resCreate, resWait *pb.Result) {
	p.log().Debug("Creating disk snapshot...")
	var err error
	if err = resCreate.Timeit(p.snapshot.create, ctx); err == nil {
		if err = resWait.Timeit(p.snapshot.createWait, ctx); err != nil {
			_ = p.snapshot.delete(ctx)
		}
	}
	if err == nil {
		p.log().Debug("Created disk snapshot", "took_create", resCreate.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not create disk snapshot", "took_create", resCreate.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *DiskSnapshotProbe[T]) getDiskSnapshot(ctx context.Context, res *pb.Result) {
	p.log().Debug("Getting disk snapshot...")
	err := res.Timeit(p.snapshot.get, ctx)
	if err == nil {
		p.log().Debug("Got disk snapshot", "took", res.Took())
	} else {
		p.log().Error("Could not get disk snapshot", "took", res.Took(), "err", err)
	}
}

func (p *DiskSnapshotProbe[T]) deleteDiskSnapshot(ctx context.Context, resDelete, resWait *pb.Result) {
	p.log().Debug("Deleting disk snapshot...")
	var err error
	if err = resDelete.Timeit(p.snapshot.delete, ctx); err == nil {
		err = resWait.Timeit(p.snapshot.deleteWait, ctx)
	}
	if err == nil {
		p.log().Debug("Deleted disk snapshot", "took_delete", resDelete.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not delete disk snapshot", "took_delete", resDelete.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *DiskSnapshotProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "snapshot", p.snapshot.String())
}
