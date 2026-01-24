package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/spectrum"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterZoneCollector("spectrum-app", CollectSpectrumApps)
}

type SpectrumApp struct {
	Client *spectrum.AppService
	ZoneID string
}

func CollectSpectrumApps(creds *types.Credentials, zone *types.Zone) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	appPage, err := client.Spectrum.Apps.List(context.TODO(), spectrum.AppListParams{
		ZoneID: cloudflare.F(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	var allResources types.Resources
	for appPage != nil && len(appPage.Result) != 0 {
		for _, appUnion := range appPage.Result {
			// The result is a union type - try to extract as AppListResponseArray
			if appArray, ok := appUnion.(spectrum.AppListResponseArray); ok {
				for _, app := range appArray {
					displayName := app.DNS.Name
					if displayName == "" {
						displayName = app.ID
					}
					res := types.Resource{
						Removable:    SpectrumApp{Client: client.Spectrum.Apps, ZoneID: zone.ID},
						ResourceID:   app.ID,
						ResourceName: displayName,
						AccountID:    creds.AccountID,
						ProductName:  "SpectrumApp",
					}
					allResources = append(allResources, &res)
				}
			}
		}
		appPage, err = appPage.GetNextPage()
		if err != nil {
			break
		}
	}

	return allResources, nil
}

func (c SpectrumApp) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, spectrum.AppDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
