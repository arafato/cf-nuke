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
	infrastructure.RegisterCollector("gateway-rule", CollectGatewayRules)
}

type GatewayRule struct {
	Client *zero_trust.GatewayRuleService
}

func CollectGatewayRules(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Gateway.Rules.List(context.TODO(), zero_trust.GatewayRuleListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		if utils.IsSkippableError(err) {
			utils.AddWarning("GatewayRule", "", "insufficient permissions or feature not available")
			return nil, nil
		}
		return nil, err
	}

	var allRules []zero_trust.GatewayRule
	for page != nil && len(page.Result) != 0 {
		allRules = append(allRules, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, rule := range allRules {
		// Skip deleted rules
		if rule.DeletedAt.IsZero() == false {
			continue
		}

		displayName := rule.Name
		if displayName == "" {
			displayName = rule.ID
		}
		res := types.Resource{
			Removable:    GatewayRule{Client: client.ZeroTrust.Gateway.Rules},
			ResourceID:   rule.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "GatewayRule",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c GatewayRule) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.GatewayRuleDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
