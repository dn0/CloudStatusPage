package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"cspage/pkg/mon/agent"
)

func NewCredential() *azidentity.DefaultAzureCredential {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		agent.Die("Could not initialize Azure default credential", "err", err)
	}
	return cred
}
