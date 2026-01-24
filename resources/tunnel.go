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
	infrastructure.RegisterAccountCollector("tunnel", CollectTunnels)
}

type Tunnel struct {
	Client *zero_trust.TunnelCloudflaredService
}

func CollectTunnels(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Tunnels.Cloudflared.List(context.TODO(), zero_trust.TunnelCloudflaredListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})
	if err != nil {
		return nil, err
	}

	var allTunnels []zero_trust.TunnelCloudflaredListResponse
	for page != nil && len(page.Result) != 0 {
		allTunnels = append(allTunnels, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, tunnel := range allTunnels {
		// Skip deleted tunnels
		if !tunnel.DeletedAt.IsZero() {
			continue
		}

		displayName := tunnel.Name
		if displayName == "" {
			displayName = tunnel.ID
		}
		res := types.Resource{
			Removable:    Tunnel{Client: client.ZeroTrust.Tunnels.Cloudflared},
			ResourceID:   tunnel.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "Tunnel",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Tunnel) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.TunnelCloudflaredDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
