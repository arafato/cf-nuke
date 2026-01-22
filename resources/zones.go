package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("zone", CollectZones)
}

type Zone struct {
	Client *zones.ZoneService
}

func CollectZones(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Zones.List(context.TODO(), zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.F(creds.AccountID)}),
	})

	var allZones []zones.Zone

	if err != nil {
		return nil, err
	}

	for page != nil && len(page.Result) != 0 {
		allZones = append(allZones, page.Result...)
		page, err = page.GetNextPage()
	}

	var allResources types.Resources
	for _, zone := range allZones {
		res := types.Resource{
			Removable:    Zone{Client: client.Zones},
			ResourceID:   zone.ID,
			ResourceName: zone.Name,
			AccountID:    creds.AccountID,
			ProductName:  "zone",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Zone) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), zones.ZoneDeleteParams{
		ZoneID: cloudflare.F(resourceID)})

	return err
}
