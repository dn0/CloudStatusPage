package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
)

type vhdSnapshot struct {
	namePrefix    string
	resourceGroup string
	client        *armcompute.SnapshotsClient
	timeout       time.Duration
	waitTimeout   time.Duration

	name   string
	exists bool

	createParams *armcompute.Snapshot
	createPoller *runtime.Poller[armcompute.SnapshotsClientCreateOrUpdateResponse]
	deletePoller *runtime.Poller[armcompute.SnapshotsClientDeleteResponse]
}

func (s *vhdSnapshot) String() string {
	return "name=" + s.name
}

func (s *vhdSnapshot) init(region, diskID string) {
	s.createParams = &armcompute.Snapshot{
		Location: to.Ptr(region),
		SKU: &armcompute.SnapshotSKU{
			Name: to.Ptr(armcompute.SnapshotStorageAccountTypesStandardLRS),
		},
		Properties: &armcompute.SnapshotProperties{
			CreationData: &armcompute.CreationData{
				CreateOption:     to.Ptr(armcompute.DiskCreateOptionCopy),
				SourceResourceID: to.Ptr(diskID),
			},
		},
		Tags: commonTags,
	}
}

func (s *vhdSnapshot) new() {
	s.name = fmt.Sprintf("%s-%d", s.namePrefix, time.Now().Unix())
	s.exists = false
}

func (s *vhdSnapshot) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var err error
	if s.createPoller, err = s.client.BeginCreateOrUpdate(ctx, s.resourceGroup, s.name, *(s.createParams), nil); err != nil {
		return fmt.Errorf("snapshot(%s).BeginCreate: %w", s, err)
	}

	s.exists = true
	return nil
}

func (s *vhdSnapshot) createWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.waitTimeout)
	defer cancel()

	if _, err := s.createPoller.PollUntilDone(ctx, pollUntilDoneOptions); err != nil {
		return fmt.Errorf("snapshot(%s).BeginCreate.PollUntilDone: %w", s, err)
	}

	s.createPoller = nil
	return nil
}

func (s *vhdSnapshot) get(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.Get(ctx, s.resourceGroup, s.name, nil); err != nil {
		return fmt.Errorf("snapshot(%s).Get: %w", s, err)
	}

	return nil
}

func (s *vhdSnapshot) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var err error
	if s.deletePoller, err = s.client.BeginDelete(ctx, s.resourceGroup, s.name, nil); err != nil {
		return fmt.Errorf("snapshot(%s).Delete: %w", s, err)
	}

	s.exists = false
	return nil
}

func (s *vhdSnapshot) deleteWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.waitTimeout)
	defer cancel()

	if _, err := s.deletePoller.PollUntilDone(ctx, pollUntilDoneOptions); err != nil {
		return fmt.Errorf("snapshot(%s).BeginDelete.PollUntilDone: %w", s, err)
	}

	s.deletePoller = nil
	return nil
}
