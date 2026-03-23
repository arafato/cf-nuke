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
	infrastructure.RegisterAccountCollector("zt-identity-provider", CollectZTIdentityProviders)
}

type ZTIdentityProvider struct {
	Client *zero_trust.IdentityProviderService
}

func CollectZTIdentityProviders(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.IdentityProviders.List(context.TODO(), zero_trust.IdentityProviderListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})
	if err != nil {
		return nil, err
	}

	var allIdPs []zero_trust.IdentityProviderListResponse
	for page != nil && len(page.Result) != 0 {
		allIdPs = append(allIdPs, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, idp := range allIdPs {
		displayName := idp.Name
		if displayName == "" {
			displayName = idp.ID
		}
		res := types.Resource{
			Removable:    ZTIdentityProvider{Client: client.ZeroTrust.IdentityProviders},
			ResourceID:   idp.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "ZTIdentityProvider",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c ZTIdentityProvider) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.IdentityProviderDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
