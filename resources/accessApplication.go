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
	infrastructure.RegisterAccountCollector("zt-access-application", CollectZTAccessApplications)
}

type ZTAccessApplication struct {
	Client *zero_trust.AccessApplicationService
}

func CollectZTAccessApplications(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	appPage, err := client.ZeroTrust.Access.Applications.List(context.TODO(), zero_trust.AccessApplicationListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})
	if err != nil {
		return nil, err
	}

	var allApps []zero_trust.AccessApplicationListResponse
	for appPage != nil && len(appPage.Result) != 0 {
		allApps = append(allApps, appPage.Result...)
		appPage, err = appPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, app := range allApps {
		// Use name as display, fall back to domain or ID
		name := app.Name
		if name == "" {
			name = app.Domain
		}
		if name == "" {
			name = app.ID
		}
		res := types.Resource{
			Removable:    ZTAccessApplication{Client: client.ZeroTrust.Access.Applications},
			ResourceID:   app.ID,
			ResourceName: name,
			AccountID:    creds.AccountID,
			ProductName:  "ZTAccessApplication",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c ZTAccessApplication) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.AccessApplicationDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
