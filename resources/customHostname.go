package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/custom_hostnames"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterZoneCollector("custom-hostname", CollectCustomHostnames)
}

type CustomHostname struct {
	Client *custom_hostnames.CustomHostnameService
	ZoneID string
}

func CollectCustomHostnames(creds *types.Credentials, zone *types.Zone) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	chPage, err := client.CustomHostnames.List(context.TODO(), custom_hostnames.CustomHostnameListParams{
		ZoneID: cloudflare.F(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	var allCustomHostnames []custom_hostnames.CustomHostnameListResponse
	for chPage != nil && len(chPage.Result) != 0 {
		allCustomHostnames = append(allCustomHostnames, chPage.Result...)
		chPage, err = chPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, ch := range allCustomHostnames {
		res := types.Resource{
			Removable:    CustomHostname{Client: client.CustomHostnames, ZoneID: zone.ID},
			ResourceID:   ch.ID,
			ResourceName: ch.Hostname,
			AccountID:    creds.AccountID,
			ProductName:  "CustomHostname",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c CustomHostname) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, custom_hostnames.CustomHostnameDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
