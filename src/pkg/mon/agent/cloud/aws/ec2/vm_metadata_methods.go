package ec2

import (
	"context"
	"fmt"

	"cspage/pkg/http"
)

const (
	vmMetadataBaseURL = "http://169.254.169.254/latest/meta-data/"
)

//nolint:gochecknoglobals // This is a constant.
var vmMetadataBaseHeaders = map[string]string{}

func getInstanceMetadata(ctx context.Context, client *http.Client, suffix string) (string, error) {
	url := vmMetadataBaseURL + suffix
	res, err := client.GetString(ctx, url, vmMetadataBaseHeaders)
	if err != nil {
		return "", fmt.Errorf("metadata(url=%s).Get(): %w", url, err)
	}
	return res, nil
}
