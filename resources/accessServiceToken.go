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
	infrastructure.RegisterAccountCollector("zt-service-token", CollectZTServiceTokens)
}

type ZTServiceToken struct {
	Client *zero_trust.AccessServiceTokenService
}

func CollectZTServiceTokens(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Access.ServiceTokens.List(context.TODO(), zero_trust.AccessServiceTokenListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allTokens []zero_trust.ServiceToken
	for page != nil && len(page.Result) != 0 {
		allTokens = append(allTokens, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, token := range allTokens {
		displayName := token.Name
		if displayName == "" {
			displayName = token.ID
		}
		res := types.Resource{
			Removable:    ZTServiceToken{Client: client.ZeroTrust.Access.ServiceTokens},
			ResourceID:   token.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "ZTServiceToken",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c ZTServiceToken) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.AccessServiceTokenDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
