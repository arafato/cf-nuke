package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/rulesets"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterAccountCollector("ruleset", CollectRulesets)
}

type Ruleset struct {
	Client *rulesets.RulesetService
}

func CollectRulesets(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Rulesets.List(context.TODO(), rulesets.RulesetListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allRulesets []rulesets.RulesetListResponse
	for page != nil && len(page.Result) != 0 {
		allRulesets = append(allRulesets, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, rs := range allRulesets {
		// Skip managed rulesets - these are Cloudflare-managed and cannot be deleted
		if rs.Kind == rulesets.KindManaged {
			continue
		}

		displayName := rs.Name
		if displayName == "" {
			displayName = rs.ID
		}
		res := types.Resource{
			Removable:    Ruleset{Client: client.Rulesets},
			ResourceID:   rs.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "Ruleset",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Ruleset) Remove(accountID string, resourceID string, resourceName string) error {
	return c.Client.Delete(context.TODO(), resourceID, rulesets.RulesetDeleteParams{
		AccountID: cloudflare.F(accountID),
	})
}
