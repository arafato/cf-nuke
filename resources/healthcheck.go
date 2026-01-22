package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/healthchecks"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("healthcheck", CollectHealthchecks)
}

type Healthcheck struct {
	Client *healthchecks.HealthcheckService
	ZoneID string
}

func CollectHealthchecks(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	// First, get all zones for the account
	zonePage, err := client.Zones.List(context.TODO(), zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.F(creds.AccountID)}),
	})
	if err != nil {
		return nil, err
	}

	var allZones []zones.Zone
	for zonePage != nil && len(zonePage.Result) != 0 {
		allZones = append(allZones, zonePage.Result...)
		zonePage, err = zonePage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources

	// For each zone, list all healthchecks
	for _, zone := range allZones {
		hcPage, err := client.Healthchecks.List(context.TODO(), healthchecks.HealthcheckListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			if utils.IsSkippableError(err) {
				utils.AddWarning("Healthcheck", zone.Name, "insufficient permissions")
			}
			continue
		}

		var allHealthchecks []healthchecks.Healthcheck
		for hcPage != nil && len(hcPage.Result) != 0 {
			allHealthchecks = append(allHealthchecks, hcPage.Result...)
			hcPage, err = hcPage.GetNextPage()
			if err != nil {
				break
			}
		}

		for _, hc := range allHealthchecks {
			displayName := hc.Name
			if displayName == "" {
				displayName = hc.Address
			}
			if displayName == "" {
				displayName = hc.ID
			}
			res := types.Resource{
				Removable:    Healthcheck{Client: client.Healthchecks, ZoneID: zone.ID},
				ResourceID:   hc.ID,
				ResourceName: displayName,
				AccountID:    creds.AccountID,
				ProductName:  "Healthcheck",
			}
			allResources = append(allResources, &res)
		}
	}

	return allResources, nil
}

func (c Healthcheck) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, healthchecks.HealthcheckDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
