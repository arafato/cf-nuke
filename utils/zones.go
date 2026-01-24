package utils

import (
	"context"
	"slices"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/types"
)

// FetchZones retrieves all zones for the account, excluding any zones in the excludeZones list
func FetchZones(creds *types.Credentials, excludeZones []string) ([]*types.Zone, error) {
	client := CreateCFClient(creds)

	page, err := client.Zones.List(context.TODO(), zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.F(creds.AccountID)}),
	})
	if err != nil {
		return nil, err
	}

	var allZones []*types.Zone
	for page != nil && len(page.Result) != 0 {
		for _, z := range page.Result {
			// Skip excluded zones
			if slices.Contains(excludeZones, z.Name) {
				continue
			}
			allZones = append(allZones, &types.Zone{
				ID:   z.ID,
				Name: z.Name,
			})
		}
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	return allZones, nil
}
