package ec2

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	ebsSnapshotProbeName             = "aws_ec2_ebs_snapshot" // doc="Amazon EBS Volume Snapshot"
	ebsSnapshotProbeActionCreate     = 10                     // name="compute.ebs.snapshot.create"      doc="Creates a new EBS snapshot" url="https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CreateSnapshot.html"
	ebsSnapshotProbeActionCreateWait = 11                     // name="compute.ebs.snapshot.createWait"  doc="Waits for a new EBS snapshot created by `compute.ebs.snapshot.create` to complete" url=""
	ebsSnapshotProbeActionGet        = 20                     // name="compute.ebs.snapshot.get"         doc="Describes the EBS snapshot created by `compute.ebs.snapshot.create`" url="https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSnapshots.html"
	ebsSnapshotProbeActionDelete     = 30                     // name="compute.ebs.snapshot.delete"      doc="Deletes the EBS snapshot created by `compute.ebs.snapshot.create`" url="https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DeleteSnapshot.html"
	// ebsSnapshotProbeActionDeleteWait = 31                  // name="compute.ebs.snapshot.deleteWait"  doc="Waits for the EBS snapshot deleted by `compute.ebs.snapshot.delete` to get completely removed" url="".
)

type EBSSnapshotProbe[T agent.AWS] struct {
	cfg      *agent.AWSConfig
	snapshot *ebsSnapshot
}

func NewEBSSnapshotProbe[T agent.AWS](cfg *agent.AWSConfig, awsConfig *aws.Config) *EBSSnapshotProbe[T] {
	return &EBSSnapshotProbe[T]{
		cfg: cfg,
		snapshot: &ebsSnapshot{
			namePrefix:  cfg.Cloud.EC2EBSSnapshotPrefix,
			client:      ec2.NewFromConfig(*awsConfig),
			timeout:     cfg.ProbeTimeout,
			waitTimeout: cfg.ProbeLongTimeout,
		},
	}
}

func (p *EBSSnapshotProbe[T]) String() string {
	return ebsSnapshotProbeName
}

func (p *EBSSnapshotProbe[T]) Start(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, p.cfg.ProbeTimeout)
	defer cancel()

	if p.cfg.Cloud.EC2EBSSnapshotVolumeID == "" {
		var err error
		if p.snapshot.volumeID, err = getVolumeID(ctx, p.snapshot.client, p.cfg.Env.VM.ID); err != nil {
			agent.DieLog(p.log(), "Could not fetch EBS volume ID", "instance_id", p.cfg.Env.VM.ID, "err", err)
		}
	} else {
		p.snapshot.volumeID = p.cfg.Cloud.EC2EBSSnapshotVolumeID
	}

	p.log().Info(
		"Probe initialized",
		"volume_id", p.snapshot.volumeID,
		"snapshot_prefix", p.snapshot.namePrefix,
	)
}

func (p *EBSSnapshotProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(ebsSnapshotProbeActionCreate),
		pb.NewResult(ebsSnapshotProbeActionCreateWait),
		pb.NewResult(ebsSnapshotProbeActionGet),
		pb.NewResult(ebsSnapshotProbeActionDelete),
	}
	p.snapshot.new()

	p.createEBSSnapshot(ctx, res[0], res[1])
	if res[0].Failed() || res[1].Failed() {
		p.cleanup(ctx)
		return res
	}
	p.getEBSSnapshot(ctx, res[2])
	p.deleteEBSSnapshot(ctx, res[3])

	return res
}

func (p *EBSSnapshotProbe[T]) Stop(ctx context.Context) {
	p.cleanup(ctx)
}

func (p *EBSSnapshotProbe[T]) cleanup(ctx context.Context) {
	if p.snapshot.id != "" {
		p.log().Debug("Deleting EBS snapshot")
		_ = p.snapshot.delete(ctx)
	}
}

func (p *EBSSnapshotProbe[T]) createEBSSnapshot(ctx context.Context, resCreate, resWait *pb.Result) {
	p.log().Debug("Creating EBS snapshot...")
	var err error
	if err = resCreate.Timeit(p.snapshot.create, ctx); err == nil {
		if err = resWait.Timeit(p.snapshot.createWait, ctx); err != nil {
			_ = p.snapshot.delete(ctx)
		}
	}
	if err == nil {
		p.log().Debug("Created EBS snapshot", "took_create", resCreate.Took(), "took_wait", resWait.Took())
	} else {
		p.log().Error("Could not create EBS snapshot", "took_create", resCreate.Took(), "took_wait", resWait.Took(), "err", err)
	}
}

func (p *EBSSnapshotProbe[T]) getEBSSnapshot(ctx context.Context, res *pb.Result) {
	p.log().Debug("Getting EBS snapshot...")
	err := res.Timeit(p.snapshot.get, ctx)
	if err == nil {
		p.log().Debug("Got EBS snapshot", "took", res.Took())
	} else {
		p.log().Error("Could not get EBS snapshot", "took", res.Took(), "err", err)
	}
}

func (p *EBSSnapshotProbe[T]) deleteEBSSnapshot(ctx context.Context, resDelete *pb.Result) {
	p.log().Debug("Deleting EBS snapshot...", "probe", p.String(), "snapshot", p.snapshot.String())
	err := resDelete.Timeit(p.snapshot.delete, ctx)
	if err == nil {
		p.log().Debug("Deleted EBS snapshot", "took_delete", resDelete.Took(), "took_wait", -1)
	} else {
		p.log().Error("Could not delete EBS snapshot", "took_delete", resDelete.Took(), "took_wait", -1, "err", err)
	}
}

func (p *EBSSnapshotProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "snapshot", p.snapshot.String())
}
