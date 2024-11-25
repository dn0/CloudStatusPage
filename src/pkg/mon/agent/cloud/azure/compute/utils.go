package compute

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"

	"cspage/pkg/http"
)

const (
	pollFrequency = time.Second

	managedDiskMetadataSuffix = "compute/storageProfile/osDisk/managedDisk/id?api-version=2019-06-01&format=text"
)

//nolint:gochecknoglobals // This is a constant.
var pollUntilDoneOptions = &runtime.PollUntilDoneOptions{
	Frequency: pollFrequency,
}

//nolint:gochecknoglobals // This is a constant.
var commonTags = map[string]*string{
	"cost-center": to.Ptr("mon-probe"),
}

func getManagedDiskId(ctx context.Context) (string, error) {
	return getInstanceMetadata(ctx, http.NewClient(), managedDiskMetadataSuffix)
}
