package compute

import (
	"context"
	"fmt"
	"time"

	cloudCompute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

type VMProvisioningModel string

const (
	VMProvisioningModelStandard VMProvisioningModel = "STANDARD"
	VMProvisioningModelSpot     VMProvisioningModel = "SPOT"

	vmProbeStartupScript = "#!/bin/bash\n\nsleep 300;\npoweroff\n"
)

type vm struct {
	namePrefix        string
	zones             []string
	zonesSkip         []string
	provisioningModel VMProvisioningModel
	machineType       string
	client            *cloudCompute.InstancesClient
	timeout           time.Duration
	waitTimeout       time.Duration

	name     string
	zone     string
	exists   bool
	rev      int
	createOp *cloudCompute.Operation
	deleteOp *cloudCompute.Operation

	createRequest *computepb.InsertInstanceRequest
	getRequest    *computepb.GetInstanceRequest
	deleteRequest *computepb.DeleteInstanceRequest
}

func (v *vm) String() string {
	return fmt.Sprintf("name=%s zone=%s mtype=%s ptype=%s", v.name, v.zone, v.machineType, v.provisioningModel)
}

func (v *vm) init(projectID, diskImage, subnet string) {
	var provisioningModel string
	switch v.provisioningModel {
	case VMProvisioningModelStandard:
		provisioningModel = computepb.Scheduling_STANDARD.String()
	case VMProvisioningModelSpot:
		provisioningModel = computepb.Scheduling_SPOT.String()
	default:
		provisioningModel = computepb.Scheduling_UNDEFINED_PROVISIONING_MODEL.String()
	}

	v.createRequest = &computepb.InsertInstanceRequest{
		Project: projectID,
		InstanceResource: &computepb.Instance{
			Scheduling: &computepb.Scheduling{
				ProvisioningModel: proto.String(provisioningModel),
			},
			Metadata: &computepb.Metadata{
				Items: []*computepb.Items{{
					Key:   proto.String("startup-script"),
					Value: proto.String(vmProbeStartupScript),
				}},
			},
			Disks: []*computepb.AttachedDisk{{
				InitializeParams: &computepb.AttachedDiskInitializeParams{
					SourceImage: proto.String(diskImage),
				},
				AutoDelete: proto.Bool(true),
				Boot:       proto.Bool(true),
				Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
			}},
			NetworkInterfaces: []*computepb.NetworkInterface{{
				Subnetwork: proto.String(subnet),
			}},
			Labels: commonLabels,
		},
	}
	v.getRequest = &computepb.GetInstanceRequest{
		Project: projectID,
	}
	v.deleteRequest = &computepb.DeleteInstanceRequest{
		Project: projectID,
	}
}

func (v *vm) new() {
	v.rev++
	v.name = fmt.Sprintf("%s-%d", v.namePrefix, time.Now().Unix())
	v.zone = v.zones[v.rev%len(v.zones)]
	v.exists = false
}

func (v *vm) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	v.createRequest.Zone = v.zone
	v.createRequest.InstanceResource.Name = proto.String(v.name)
	v.createRequest.InstanceResource.MachineType = proto.String("zones/" + v.zone + "/machineTypes/" + v.machineType)

	var err error
	if v.createOp, err = v.client.Insert(ctx, v.createRequest); err != nil {
		return fmt.Errorf("vm(%s).Insert: %w", v, err)
	}

	v.exists = true
	return nil
}

func (v *vm) createWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.waitTimeout)
	defer cancel()

	if err := opWait(ctx, v.createOp); err != nil {
		return fmt.Errorf("vm(%s).Insert.wait: %w", v, err)
	}

	v.createOp = nil
	return nil
}

func (v *vm) get(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	v.getRequest.Zone = v.zone
	v.getRequest.Instance = v.name

	if _, err := v.client.Get(ctx, v.getRequest); err != nil {
		return fmt.Errorf("vm(%s).Get: %w", v, err)
	}

	return nil
}

func (v *vm) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	v.deleteRequest.Zone = v.zone
	v.deleteRequest.Instance = v.name

	var err error
	if v.deleteOp, err = v.client.Delete(ctx, v.deleteRequest); err != nil {
		return fmt.Errorf("vm(%s).Delete: %w", v, err)
	}

	v.exists = false
	return nil
}

func (v *vm) deleteWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.waitTimeout)
	defer cancel()

	if err := opWait(ctx, v.deleteOp); err != nil {
		return fmt.Errorf("vm(%s).Delete.wait: %w", v, err)
	}

	v.deleteOp = nil
	return nil
}
