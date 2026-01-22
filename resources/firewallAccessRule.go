package resources

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/firewall"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("firewall-access-rule", CollectFirewallAccessRules)
}

type FirewallAccessRule struct {
	Client *firewall.AccessRuleService
}

func CollectFirewallAccessRules(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	// List account-level firewall access rules
	page, err := client.Firewall.AccessRules.List(context.TODO(), firewall.AccessRuleListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		if utils.IsSkippableError(err) {
			utils.AddWarning("FirewallAccessRule", "", "insufficient permissions or feature not available")
			return nil, nil
		}
		return nil, err
	}

	var allRules []firewall.AccessRuleListResponse
	for page != nil && len(page.Result) != 0 {
		allRules = append(allRules, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, rule := range allRules {
		// Create descriptive name from mode and configuration
		displayName := rule.Notes
		if displayName == "" {
			displayName = fmt.Sprintf("%s %s", rule.Mode, rule.Configuration.Value)
		}
		res := types.Resource{
			Removable:    FirewallAccessRule{Client: client.Firewall.AccessRules},
			ResourceID:   rule.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "FirewallAccessRule",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c FirewallAccessRule) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, firewall.AccessRuleDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
