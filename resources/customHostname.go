package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/custom_hostnames"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("custom-hostname", CollectCustomHostnames)
}

type CustomHostname struct {
	Client *custom_hostnames.CustomHostnameService
	ZoneID string
}

func CollectCustomHostnames(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	// First, get all zones for the account
	zonePage, err := client.Zones.List(context.TODO(), zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.F(creds.AccountID)}),
	})
	if err != nil {
		return nil, err
	}

	var allZones []zones.Zone
	for len(zonePage.Result) != 0 {
		allZones = append(allZones, zonePage.Result...)
		zonePage, err = zonePage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources

	// For each zone, list all custom hostnames
	for _, zone := range allZones {
		chPage, err := client.CustomHostnames.List(context.TODO(), custom_hostnames.CustomHostnameListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			// Skip zones where we might not have permissions
			continue
		}

		var allCustomHostnames []custom_hostnames.CustomHostnameListResponse
		for len(chPage.Result) != 0 {
			allCustomHostnames = append(allCustomHostnames, chPage.Result...)
			chPage, err = chPage.GetNextPage()
			if err != nil {
				break
			}
		}

		for _, ch := range allCustomHostnames {
			res := types.Resource{
				Removable:    CustomHostname{Client: client.CustomHostnames, ZoneID: zone.ID},
				ResourceID:   ch.ID,
				ResourceName: ch.Hostname,
				AccountID:    creds.AccountID,
				ProductName:  "CustomHostname",
				State:        types.Ready,
			}
			allResources = append(allResources, &res)
		}
	}

	return allResources, nil
}

func (c CustomHostname) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, custom_hostnames.CustomHostnameDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
