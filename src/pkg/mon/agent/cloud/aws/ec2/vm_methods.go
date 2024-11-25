package ec2

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type VMProvisioningModel string

const (
	VMProvisioningModelStandard VMProvisioningModel = "STANDARD"
	VMProvisioningModelSpot     VMProvisioningModel = "SPOT"

	vmProbeNewInstanceId = "TBD"
	vmProbeStartupScript = "IyEvYmluL2Jhc2gKCnNsZWVwIDMwMDsKcG93ZXJvZmYK" // sleep 300; poweroff
)

type vm struct {
	namePrefix        string
	zones             []string
	zonesSkip         []string
	provisioningModel VMProvisioningModel
	machineType       string
	imageName         string
	imageID           string
	subnets           map[string]string
	client            *ec2.Client
	timeout           time.Duration
	waitTimeout       time.Duration

	id     string
	name   string
	zone   string
	exists bool
	rev    int

	createInput *ec2.RunInstancesInput
}

func (v *vm) String() string {
	return fmt.Sprintf("id=%s zone=%s mtype=%s, ptype=%s", v.id, v.zone, v.machineType, v.provisioningModel)
}

func (v *vm) new() {
	v.rev++
	v.id = vmProbeNewInstanceId
	v.name = fmt.Sprintf("%s-%d", v.namePrefix, time.Now().Unix())
	v.zone = v.zones[v.rev%len(v.zones)]
	v.exists = false
	// New input has to be created for every RunInstances call
	v.createInput = &ec2.RunInstancesInput{
		Placement: &types.Placement{
			AvailabilityZone: aws.String(v.zone),
		},
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags:         resourceTags(v.name),
			},
			{
				ResourceType: types.ResourceTypeVolume,
				Tags:         resourceTags(v.name),
			},
		},
		InstanceType:                      types.InstanceType(v.machineType),
		MinCount:                          aws.Int32(1),
		MaxCount:                          aws.Int32(1),
		ImageId:                           aws.String(v.imageID),
		UserData:                          aws.String(vmProbeStartupScript),
		InstanceInitiatedShutdownBehavior: types.ShutdownBehaviorTerminate,
		PrivateDnsNameOptions: &types.PrivateDnsNameOptionsRequest{
			EnableResourceNameDnsARecord: aws.Bool(true),
			HostnameType:                 types.HostnameTypeResourceName,
		},
		NetworkInterfaces: []types.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:         aws.Int32(0),
				SubnetId:            aws.String(v.subnets[v.zone]),
				DeleteOnTermination: aws.Bool(true),
			},
		},
	}

	if v.provisioningModel == VMProvisioningModelSpot {
		v.createInput.InstanceMarketOptions = &types.InstanceMarketOptionsRequest{
			MarketType: types.MarketTypeSpot,
		}
	}
}

func (v *vm) create(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	res, err := v.client.RunInstances(ctx, v.createInput)
	if err != nil {
		// If this fails on context.Cancel/DeadlineExceeded then there is a chance that the instance was created
		return fmt.Errorf("vm(%s).Run: %w", v, err)
	}

	v.id = aws.ToString(res.Instances[0].InstanceId)
	v.exists = true
	return nil
}

func (v *vm) createWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.waitTimeout)
	defer cancel()

	waiter := ec2.NewInstanceRunningWaiter(v.client, setRunningWaiterOptions)
	if err := waiter.Wait(ctx, &ec2.DescribeInstancesInput{InstanceIds: []string{v.id}}, v.waitTimeout); err != nil {
		return fmt.Errorf("vm(%s).Run.wait: %w", v, err)
	}

	return nil
}

func (v *vm) get(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	if _, err := v.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{InstanceIds: []string{v.id}}); err != nil {
		return fmt.Errorf("vm(%s).Describe: %w", v, err)
	}

	return nil
}

func (v *vm) delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	if _, err := v.client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{InstanceIds: []string{v.id}}); err != nil {
		return fmt.Errorf("vm(%s).Terminate: %w", v, err)
	}

	v.exists = false
	return nil
}

func (v *vm) deleteWait(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.waitTimeout)
	defer cancel()

	waiter := ec2.NewInstanceTerminatedWaiter(v.client, setTerminatedWaiterOptions)
	if err := waiter.Wait(ctx, &ec2.DescribeInstancesInput{InstanceIds: []string{v.id}}, v.waitTimeout); err != nil {
		return fmt.Errorf("vm(%s).Terminate.wait: %w", v, err)
	}

	return nil
}

func setRunningWaiterOptions(options *ec2.InstanceRunningWaiterOptions) {
	options.MinDelay = waiterMinDelay
	options.MaxDelay = waiterMaxDelay
}

func setTerminatedWaiterOptions(options *ec2.InstanceTerminatedWaiterOptions) {
	options.MinDelay = waiterMinDelay
	options.MaxDelay = waiterMaxDelay
}
