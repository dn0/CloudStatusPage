package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
)

type VMProvisioningModel string

const (
	VMProvisioningModelStandard VMProvisioningModel = "STANDARD"
	VMProvisioningModelSpot     VMProvisioningModel = "SPOT"

	vmProbeStartupScript = "IyEvYmluL2Jhc2gKCnNsZWVwIDMwMDsKcG93ZXJvZmYK" // sleep 300; poweroff
)

type vm struct {
	namePrefix        string
	resourceGroup     string
	provisioningModel VMProvisioningModel
	machineType       string
	nicID             string
	client            *armcompute.VirtualMachinesClient
	timeout           time.Duration
	waitTimeout       time.Duration

	name   string
	exists bool

	createParams  *armcompute.VirtualMachine
	createPoller  *runtime.Poller[armcompute.VirtualMachinesClientCreateOrUpdateResponse]
	deletePoller  *runtime.Poller[armcompute.VirtualMachinesClientDeleteResponse]
	deleteOptions *armcompute.VirtualMachinesClientBeginDeleteOptions
}

func (v *vm) String() string {
	return fmt.Sprintf("name=%s zone=- mtype=%s, ptype=%s", v.name, v.machineType, v.provisioningModel)
}

//nolint:funlen // Create VM requires this many things :(
func (v *vm) init(region, imgPublisher, imgOffer, imgSKU, imgVersion string) {
	var vmPriority armcompute.VirtualMachinePriorityTypes
	switch v.provisioningModel {
	case VMProvisioningModelStandard:
		vmPriority = armcompute.VirtualMachinePriorityTypesRegular
	case VMProvisioningModelSpot:
		vmPriority = armcompute.VirtualMachinePriorityTypesSpot
	}

	//goland:noinspection SpellCheckingInspection
	v.createParams = &armcompute.VirtualMachine{
		Location: to.Ptr(region),
		Identity: &armcompute.VirtualMachineIdentity{
			Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
		},
		Properties: &armcompute.VirtualMachineProperties{
			Priority: to.Ptr(vmPriority),
			UserData: to.Ptr(vmProbeStartupScript),
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(armcompute.VirtualMachineSizeTypes(v.machineType)),
			},
			OSProfile: &armcompute.OSProfile{
				ComputerName:  to.Ptr(v.namePrefix + "-vm"),
				AdminUsername: to.Ptr("azureuser"),
				LinuxConfiguration: &armcompute.LinuxConfiguration{
					DisablePasswordAuthentication: to.Ptr(true),
					SSH: &armcompute.SSHConfiguration{
						PublicKeys: []*armcompute.SSHPublicKey{{
							Path: to.Ptr("/home/azureuser/.ssh/authorized_keys"),
							//nolint:lll // Just some random SSH key which nobody has access to.
							KeyData: to.Ptr("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC/AF/+H/N/97Rb98Gc5MexJqm3WjHpErMJwrmY3m36BR/LU3bfYGMAaJDZmyeM7A0cRl4beM1Cs/EWGcD0RNSlJ6sydTrnwSyjAK0LKuuAGefjNQfxAqjGgY6oDRzcPvfi5YZkvp0SpENKX7C+PKT4K8WZxMuvZp9msutAmEXvmWisNQ83kfpJ0nNrFnuLdniZQ6VXEGL2njtIOBrOYeFgdOVq681MxMBDInPB3/IrFkBd0DeHa6U0wUsRRcM+5EpoAEAGDr/Q7hRyugUxToHV/LyGCxF7UOkzuaHdIZ1+QAJ9kj7XGTGK4XSXUhOanOzL0WmNFBjAOPHJ3DCoYkAX"),
						}},
					},
				},
			},
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					Publisher: to.Ptr(imgPublisher),
					Offer:     to.Ptr(imgOffer),
					SKU:       to.Ptr(imgSKU),
					Version:   to.Ptr(imgVersion),
				},
				OSDisk: &armcompute.OSDisk{
					Name:         to.Ptr(v.namePrefix + "-disk"),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					DeleteOption: to.Ptr(armcompute.DiskDeleteOptionTypesDelete),
				},
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{{
					ID: to.Ptr(v.nicID),
					Properties: &armcompute.NetworkInterfaceReferenceProperties{
						Primary:      to.Ptr(true),
						DeleteOption: to.Ptr(armcompute.DeleteOptionsDetach),
					},
				}},
			},
		},
		Tags: commonTags,
	}
	v.deleteOptions = &armcompute.VirtualMachinesClientBeginDeleteOptions{
		ForceDeletion: to.Ptr(true),
	}
}

func (v *vm) new() {
	v.name = fmt.Sprintf("%s-%d", v.namePrefix, time.Now().Unix())
	v.exists = false
}

func (v *vm) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	v.createParams.Properties.OSProfile.ComputerName = to.Ptr(v.name)

	var err error
	if v.createPoller, err = v.client.BeginCreateOrUpdate(ctx, v.resourceGroup, v.name, *(v.createParams), nil); err != nil {
		return fmt.Errorf("vm(%s).BeginCreate: %w", v, err)
	}

	v.exists = true
	return nil
}

func (v *vm) createWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.waitTimeout)
	defer cancel()

	if _, err := v.createPoller.PollUntilDone(ctx, pollUntilDoneOptions); err != nil {
		return fmt.Errorf("vm(%s).BeginCreate.PollUntilDone: %w", v, err)
	}

	v.createPoller = nil
	return nil
}

func (v *vm) get(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	if _, err := v.client.Get(ctx, v.resourceGroup, v.name, nil); err != nil {
		return fmt.Errorf("vm(%s).Get: %w", v, err)
	}

	return nil
}

func (v *vm) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	var err error
	if v.deletePoller, err = v.client.BeginDelete(ctx, v.resourceGroup, v.name, v.deleteOptions); err != nil {
		return fmt.Errorf("vm(%s).BeginDelete: %w", v, err)
	}

	v.exists = false
	return nil
}

func (v *vm) deleteWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.waitTimeout)
	defer cancel()

	if _, err := v.deletePoller.PollUntilDone(ctx, pollUntilDoneOptions); err != nil {
		return fmt.Errorf("vm(%s).BeginDelete.PollUntilDone: %w", v, err)
	}

	v.deletePoller = nil
	return nil
}
