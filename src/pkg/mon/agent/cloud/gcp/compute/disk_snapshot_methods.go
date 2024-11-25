package compute

import (
	"context"
	"fmt"
	"time"

	cloudCompute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

type diskSnapshot struct {
	namePrefix  string
	client      *cloudCompute.SnapshotsClient
	timeout     time.Duration
	waitTimeout time.Duration

	name     string
	exists   bool
	createOp *cloudCompute.Operation
	deleteOp *cloudCompute.Operation

	createRequest *computepb.InsertSnapshotRequest
	getRequest    *computepb.GetSnapshotRequest
	deleteRequest *computepb.DeleteSnapshotRequest
}

func (s *diskSnapshot) String() string {
	return "name=" + s.name
}

func (s *diskSnapshot) init(projectID, region, zone, diskName string) {
	s.createRequest = &computepb.InsertSnapshotRequest{
		Project: projectID,
		SnapshotResource: &computepb.Snapshot{
			SourceDisk:       proto.String("projects/" + projectID + "/zones/" + zone + "/disks/" + diskName),
			StorageLocations: []string{region},
			Labels:           commonLabels,
		},
	}
	s.getRequest = &computepb.GetSnapshotRequest{
		Project: projectID,
	}
	s.deleteRequest = &computepb.DeleteSnapshotRequest{
		Project: projectID,
	}
}

func (s *diskSnapshot) new() {
	s.name = fmt.Sprintf("%s-%d", s.namePrefix, time.Now().Unix())
	s.exists = false
}

func (s *diskSnapshot) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	s.createRequest.SnapshotResource.Name = proto.String(s.name)

	var err error
	if s.createOp, err = s.client.Insert(ctx, s.createRequest); err != nil {
		return fmt.Errorf("snapshot(%s).Insert: %w", s, err)
	}

	s.exists = true
	return nil
}

func (s *diskSnapshot) createWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.waitTimeout)
	defer cancel()

	if err := opWait(ctx, s.createOp); err != nil {
		return fmt.Errorf("snapshot(%s).Insert.wait: %w", s, err)
	}

	s.createOp = nil
	return nil
}

func (s *diskSnapshot) get(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	s.getRequest.Snapshot = s.name

	if _, err := s.client.Get(ctx, s.getRequest); err != nil {
		return fmt.Errorf("snapshot(%s).Get: %w", s, err)
	}

	return nil
}

func (s *diskSnapshot) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	s.deleteRequest.Snapshot = s.name

	var err error
	if s.deleteOp, err = s.client.Delete(ctx, s.deleteRequest); err != nil {
		return fmt.Errorf("snapshot(%s).Delete: %w", s, err)
	}

	s.exists = false
	return nil
}

func (s *diskSnapshot) deleteWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.waitTimeout)
	defer cancel()

	if err := opWait(ctx, s.deleteOp); err != nil {
		return fmt.Errorf("snapshot(%s).Delete.wait: %w", s, err)
	}

	s.deleteOp = nil
	return nil
}
