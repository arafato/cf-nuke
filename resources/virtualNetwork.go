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
	infrastructure.RegisterCollector("virtual-network", CollectVirtualNetworks)
}

type VirtualNetwork struct {
	Client *zero_trust.NetworkVirtualNetworkService
}

func CollectVirtualNetworks(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Networks.VirtualNetworks.List(context.TODO(), zero_trust.NetworkVirtualNetworkListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		if utils.IsSkippableError(err) {
			utils.AddWarning("VirtualNetwork", "", "insufficient permissions or feature not available")
			return nil, nil
		}
		return nil, err
	}

	var allVNets []zero_trust.VirtualNetwork
	for page != nil && len(page.Result) != 0 {
		allVNets = append(allVNets, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, vnet := range allVNets {
		// Skip deleted virtual networks
		if !vnet.DeletedAt.IsZero() {
			continue
		}

		// Skip default network - it can't be deleted
		if vnet.IsDefaultNetwork {
			continue
		}

		displayName := vnet.Name
		if displayName == "" {
			displayName = vnet.ID
		}
		res := types.Resource{
			Removable:    VirtualNetwork{Client: client.ZeroTrust.Networks.VirtualNetworks},
			ResourceID:   vnet.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "VirtualNetwork",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c VirtualNetwork) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.NetworkVirtualNetworkDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
