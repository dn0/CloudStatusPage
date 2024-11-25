package ec2

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const (
	waiterMinDelay = 100 * time.Millisecond
	waiterMaxDelay = 1000 * time.Millisecond
)

var (
	errNoAvailabilityZones  = errors.New("no availability zones found")
	errInstanceNotFound     = errors.New("AWS instance not found")
	errMultipleSubnetsFound = errors.New("multiple subnets in one availability zone found")
)

//nolint:gochecknoglobals // This is a constant.
var commonTags = []types.Tag{
	{
		Key:   aws.String("cost-center"),
		Value: aws.String("mon-probe"),
	},
	{
		Key:   aws.String("owner"),
		Value: aws.String("mon-agent"),
	},
}

func resourceTags(name string) []types.Tag {
	return append(commonTags, types.Tag{
		Key:   aws.String("Name"),
		Value: aws.String(name),
	})
}

func getLatestAMI(ctx context.Context, client *ec2.Client, filters ...types.Filter) (string, error) {
	res, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{Filters: filters})
	if err != nil {
		return "", fmt.Errorf("ec2.Client.DescribeImages: %w", err)
	}

	var errr error
	//nolint:varnamelen // Variable names a & b make sense in this context.
	slices.SortFunc(res.Images, func(a, b types.Image) int {
		ta, terr := time.Parse(time.RFC3339, aws.ToString(a.CreationDate))
		if terr != nil {
			errr = terr
			return 0
		}
		tb, terr := time.Parse(time.RFC3339, aws.ToString(b.CreationDate))
		if terr != nil {
			errr = terr
			return 0
		}
		return tb.Compare(ta)
	})

	return aws.ToString(res.Images[0].ImageId), errr
}

func getAvailabilityZones(ctx context.Context, client *ec2.Client, skip []string) ([]string, error) {
	res, err := client.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{
		AllAvailabilityZones: aws.Bool(false),
	})
	if err != nil {
		return nil, fmt.Errorf("ec2.Client.DescribeAvailabilityZones: %w", err)
	}

	var zones []string
	//nolint:gocritic // This code is used once during initialization so rangeValCopy is OK.
	for _, zone := range res.AvailabilityZones {
		if zone.State == types.AvailabilityZoneStateAvailable {
			if !slices.Contains(skip, aws.ToString(zone.ZoneName)) {
				zones = append(zones, aws.ToString(zone.ZoneName))
			}
		}
	}

	if len(zones) == 0 {
		return zones, fmt.Errorf("%s, %w", client.Options().Region, errNoAvailabilityZones)
	}

	return zones, nil
}

func getVPCSubnets(ctx context.Context, client *ec2.Client, filters ...types.Filter) (map[string]string, error) {
	res, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{Filters: filters})
	if err != nil {
		return nil, fmt.Errorf("ec2.Client.DescribeSubnets: %w", err)
	}

	subnets := make(map[string]string)
	for i := range res.Subnets {
		zone := *res.Subnets[i].AvailabilityZone
		if _, ok := subnets[zone]; ok {
			return nil, fmt.Errorf("%s: %w", zone, errMultipleSubnetsFound)
		}
		subnets[zone] = *res.Subnets[i].SubnetId
	}

	return subnets, nil
}

func getVolumeID(ctx context.Context, client *ec2.Client, instanceID string) (string, error) {
	res, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{InstanceIds: []string{instanceID}})
	if err != nil {
		return "", fmt.Errorf("ec2.Client.DescribeInstances: %w", err)
	}

	if len(res.Reservations) != 1 || len(res.Reservations[0].Instances) != 1 {
		return "", fmt.Errorf("%s: %w", instanceID, errInstanceNotFound)
	}
	return aws.ToString(res.Reservations[0].Instances[0].BlockDeviceMappings[0].Ebs.VolumeId), nil
}
