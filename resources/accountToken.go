package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/accounts"
	"github.com/cloudflare/cloudflare-go/v6/shared"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("account-token", CollectAccountToken)
}

type AccountToken struct {
	Client *accounts.AccountService
}

func CollectAccountToken(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Accounts.Tokens.List(context.TODO(), accounts.TokenListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allTokens []shared.Token

	for len(page.Result) != 0 {
		allTokens = append(allTokens, page.Result...)
		page, err = page.GetNextPage()
	}

	var allResources types.Resources
	for _, token := range allTokens {
		res := types.Resource{
			Removable:    AccountToken{Client: client.Accounts},
			ResourceID:   token.ID,
			ResourceName: token.Name,
			AccountID:    creds.AccountID,
			ProductName:  "AccountToken",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c AccountToken) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Tokens.Delete(context.TODO(), resourceID, accounts.TokenDeleteParams{
		AccountID: cloudflare.F(accountID)})

	return err
}
