package ec2

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const (
	ebsSnapshotProbeNewId = "TBD"
)

type ebsSnapshot struct {
	namePrefix  string
	volumeID    string
	client      *ec2.Client
	timeout     time.Duration
	waitTimeout time.Duration

	id     string
	name   string
	exists bool

	createInput *ec2.CreateSnapshotInput
}

func (s *ebsSnapshot) String() string {
	return fmt.Sprintf("id=%s name=%s", s.id, s.name)
}

func (s *ebsSnapshot) new() {
	s.id = ebsSnapshotProbeNewId
	s.name = fmt.Sprintf("%s-%d", s.namePrefix, time.Now().Unix())
	s.exists = false
	s.createInput = &ec2.CreateSnapshotInput{
		VolumeId: aws.String(s.volumeID),
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeSnapshot,
			Tags:         resourceTags(s.name),
		}},
	}
}

func (s *ebsSnapshot) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.CreateSnapshot(ctx, s.createInput)
	if err != nil {
		// If this fails on context.Cancel/DeadlineExceeded then there is a chance that the snapshot was created
		return fmt.Errorf("snapshot(%s).Create: %w", s, err)
	}

	s.id = aws.ToString(res.SnapshotId)
	s.exists = true
	return nil
}

func (s *ebsSnapshot) createWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.waitTimeout)
	defer cancel()

	waiter := ec2.NewSnapshotCompletedWaiter(s.client, setSnapshotCompletedWaiterOptions)
	if err := waiter.Wait(ctx, &ec2.DescribeSnapshotsInput{SnapshotIds: []string{s.id}}, s.waitTimeout); err != nil {
		return fmt.Errorf("snapshot(%s).Run.wait: %w", s, err)
	}

	return nil
}

func (s *ebsSnapshot) get(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{SnapshotIds: []string{s.id}}); err != nil {
		return fmt.Errorf("snapshot(%s).Describe: %w", s, err)
	}

	return nil
}

func (s *ebsSnapshot) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{SnapshotId: aws.String(s.id)}); err != nil {
		return fmt.Errorf("snapshot(%s).Terminate: %w", s, err)
	}

	s.exists = false
	return nil
}

func setSnapshotCompletedWaiterOptions(options *ec2.SnapshotCompletedWaiterOptions) {
	options.MinDelay = waiterMinDelay
	options.MaxDelay = waiterMaxDelay
}
