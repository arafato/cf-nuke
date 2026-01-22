package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/load_balancers"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("load-balancer", CollectLoadBalancers)
}

type LoadBalancer struct {
	Client *load_balancers.LoadBalancerService
	ZoneID string
}

func CollectLoadBalancers(creds *types.Credentials) (types.Resources, error) {
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

	// For each zone, list all load balancers
	for _, zone := range allZones {
		lbPage, err := client.LoadBalancers.List(context.TODO(), load_balancers.LoadBalancerListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			// Skip zones where we might not have permissions
			continue
		}

		var allLoadBalancers []load_balancers.LoadBalancer
		for lbPage != nil && len(lbPage.Result) != 0 {
			allLoadBalancers = append(allLoadBalancers, lbPage.Result...)
			lbPage, err = lbPage.GetNextPage()
			if err != nil {
				break
			}
		}

		for _, lb := range allLoadBalancers {
			res := types.Resource{
				Removable:    LoadBalancer{Client: client.LoadBalancers, ZoneID: zone.ID},
				ResourceID:   lb.ID,
				ResourceName: lb.Name,
				AccountID:    creds.AccountID,
				ProductName:  "LoadBalancer",
			}
			allResources = append(allResources, &res)
		}
	}

	return allResources, nil
}

func (c LoadBalancer) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, load_balancers.LoadBalancerDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
