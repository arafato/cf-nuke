package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/page_shield"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("page-shield-policy", CollectPageShieldPolicies)
}

type PageShieldPolicy struct {
	Client *page_shield.PolicyService
	ZoneID string
}

func CollectPageShieldPolicies(creds *types.Credentials) (types.Resources, error) {
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

	// For each zone, list all Page Shield policies
	for _, zone := range allZones {
		policyPage, err := client.PageShield.Policies.List(context.TODO(), page_shield.PolicyListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			if utils.IsSkippableError(err) {
				utils.AddWarning("PageShieldPolicy", zone.Name, "insufficient permissions or feature not available")
			}
			continue
		}

		var allPolicies []page_shield.PolicyListResponse
		for policyPage != nil && len(policyPage.Result) != 0 {
			allPolicies = append(allPolicies, policyPage.Result...)
			policyPage, err = policyPage.GetNextPage()
			if err != nil {
				break
			}
		}

		for _, policy := range allPolicies {
			displayName := policy.Description
			if displayName == "" {
				displayName = policy.ID
			}
			res := types.Resource{
				Removable:    PageShieldPolicy{Client: client.PageShield.Policies, ZoneID: zone.ID},
				ResourceID:   policy.ID,
				ResourceName: displayName,
				AccountID:    creds.AccountID,
				ProductName:  "PageShieldPolicy",
			}
			allResources = append(allResources, &res)
		}
	}

	return allResources, nil
}

func (c PageShieldPolicy) Remove(accountID string, resourceID string, resourceName string) error {
	return c.Client.Delete(context.TODO(), resourceID, page_shield.PolicyDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})
}
