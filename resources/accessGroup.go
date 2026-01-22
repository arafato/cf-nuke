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
	infrastructure.RegisterCollector("zt-access-group", CollectZTAccessGroups)
}

type ZTAccessGroup struct {
	Client *zero_trust.AccessGroupService
}

func CollectZTAccessGroups(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Access.Groups.List(context.TODO(), zero_trust.AccessGroupListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		if utils.IsSkippableError(err) {
			utils.AddWarning("ZTAccessGroup", "", "insufficient permissions or feature not available")
			return nil, nil
		}
		return nil, err
	}

	var allGroups []zero_trust.AccessGroupListResponse
	for page != nil && len(page.Result) != 0 {
		allGroups = append(allGroups, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, group := range allGroups {
		displayName := group.Name
		if displayName == "" {
			displayName = group.ID
		}
		res := types.Resource{
			Removable:    ZTAccessGroup{Client: client.ZeroTrust.Access.Groups},
			ResourceID:   group.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "ZTAccessGroup",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c ZTAccessGroup) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.AccessGroupDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
