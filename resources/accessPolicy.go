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
	infrastructure.RegisterAccountCollector("zt-access-policy", CollectZTAccessPolicies)
}

type ZTAccessPolicy struct {
	Client *zero_trust.AccessPolicyService
}

func CollectZTAccessPolicies(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Access.Policies.List(context.TODO(), zero_trust.AccessPolicyListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})
	if err != nil {
		return nil, err
	}

	var allPolicies []zero_trust.AccessPolicyListResponse
	for page != nil && len(page.Result) != 0 {
		allPolicies = append(allPolicies, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, policy := range allPolicies {
		displayName := policy.Name
		if displayName == "" {
			displayName = policy.ID
		}
		res := types.Resource{
			Removable:    ZTAccessPolicy{Client: client.ZeroTrust.Access.Policies},
			ResourceID:   policy.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "ZTAccessPolicy",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c ZTAccessPolicy) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.AccessPolicyDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
