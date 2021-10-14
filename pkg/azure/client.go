package azure

import (
	"fmt"
	"regexp"

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

type ARMResourceID struct {
	Subscription  string
	ResourceGroup string
	ResourceName  string
}

var regexARMResourceID = regexp.MustCompile(`(?i)^/subscriptions/(?P<subscription>[^/]*)/resourceGroups/(?P<resourceGroup>[^/]*)/providers/(?P<provider>[^/]*)/(?P<resourceType>[^/]*)/(?P<resourceName>[^/]*)$`)

// ParseARMResourceID parses a resource ID string into an ARMResourceID.
func ParseARMResourceID(id string) (*ARMResourceID, error) {
	matches := regexARMResourceID.FindAllStringSubmatch(id, -1)
	if len(matches) < 1 || len(matches[0]) < 5 {
		return nil, fmt.Errorf("parse resource id %s: not match", id)
	}

	for _, m := range matches {
		return &ARMResourceID{
			Subscription:  m[1],
			ResourceGroup: m[2],
			ResourceName:  m[5],
		}, nil
	}

	return nil, fmt.Errorf("parse resource id %s: not match", id)
}
