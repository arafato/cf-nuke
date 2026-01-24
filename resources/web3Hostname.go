package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/web3"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterZoneCollector("web3-hostname", CollectWeb3Hostnames)
}

type Web3Hostname struct {
	Client *web3.HostnameService
	ZoneID string
}

func CollectWeb3Hostnames(creds *types.Credentials, zone *types.Zone) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	hostnamePage, err := client.Web3.Hostnames.List(context.TODO(), web3.HostnameListParams{
		ZoneID: cloudflare.F(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	var allHostnames []web3.Hostname
	for hostnamePage != nil && len(hostnamePage.Result) != 0 {
		allHostnames = append(allHostnames, hostnamePage.Result...)
		hostnamePage, err = hostnamePage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, hostname := range allHostnames {
		displayName := hostname.Name
		if displayName == "" {
			displayName = hostname.ID
		}
		res := types.Resource{
			Removable:    Web3Hostname{Client: client.Web3.Hostnames, ZoneID: zone.ID},
			ResourceID:   hostname.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "Web3Hostname",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Web3Hostname) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, web3.HostnameDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
