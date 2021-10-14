package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-08-01/containerservice"
	"github.com/Azure/go-autorest/autorest"
)

func aksManagedClusterClient(
	authorizer autorest.Authorizer,
	subscriptionID string,
) containerservice.ManagedClustersClient {
	client := containerservice.NewManagedClustersClient(subscriptionID)
	client.Authorizer = authorizer

	return client
}
