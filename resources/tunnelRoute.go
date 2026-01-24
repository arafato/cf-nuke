package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/zero_trust"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterAccountCollector("tunnel-route", CollectTunnelRoutes)
}

type TunnelRoute struct {
	Client *zero_trust.NetworkRouteService
}

func CollectTunnelRoutes(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Networks.Routes.List(context.TODO(), zero_trust.NetworkRouteListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allRoutes []zero_trust.Teamnet
	for page != nil && len(page.Result) != 0 {
		allRoutes = append(allRoutes, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, route := range allRoutes {
		// Skip deleted routes
		if !route.DeletedAt.IsZero() {
			continue
		}

		displayName := route.Network
		if route.Comment != "" {
			displayName = route.Comment + " (" + route.Network + ")"
		}
		res := types.Resource{
			Removable:    TunnelRoute{Client: client.ZeroTrust.Networks.Routes},
			ResourceID:   route.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "TunnelRoute",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c TunnelRoute) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.NetworkRouteDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
