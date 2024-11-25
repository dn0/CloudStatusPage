package compute

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	cloudCompute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
)

const (
	waitMinDelay = 100 * time.Millisecond
	waitMaxDelay = 2000 * time.Millisecond
)

//nolint:gochecknoglobals // This is a constant.
var commonLabels = map[string]string{
	"cost-center": "mon-probe",
}

var errNoAvailabilityZones = errors.New("no availability zones found")

func getAvailabilityZones(ctx context.Context, projectID, region string, skip []string) ([]string, error) {
	client, err := cloudCompute.NewRegionZonesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("compute.NewRegionZonesRESTClient: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer client.Close()

	iter := client.List(ctx, &computepb.ListRegionZonesRequest{
		Project: projectID,
		Region:  region,
	})
	var zones []string
	for {
		res, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return zones, fmt.Errorf("compute.RegionZonesRESTClient.List: %w", err)
		}
		resRegion := res.GetRegion()
		// Example: https://www.googleapis.com/compute/v1/projects/cloudstatus-probe-t/regions/europe-west1
		if resRegion[strings.LastIndex(resRegion, "/")+1:] != region {
			continue
		}
		if res.GetStatus() == "UP" && res.GetDeprecated() == nil {
			if !slices.Contains(skip, res.GetName()) {
				zones = append(zones, res.GetName())
			}
		}
	}

	if len(zones) == 0 {
		return zones, fmt.Errorf("%s: %w", region, errNoAvailabilityZones)
	}

	return zones, nil
}

//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
func opWait(ctx context.Context, operation *cloudCompute.Operation) error {
	backoff := gax.Backoff{
		Initial: waitMinDelay,
		Max:     waitMaxDelay,
	}
	for {
		if err := operation.Poll(ctx); err != nil {
			return err
		}
		if operation.Done() {
			return nil
		}
		if err := gax.Sleep(ctx, backoff.Pause()); err != nil {
			return err
		}
	}
}
