package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/calls"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterAccountCollector("calls-app", CollectCallsApps)
}

type CallsApp struct {
	Client *calls.SFUService
}

func CollectCallsApps(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Calls.SFU.List(context.TODO(), calls.SFUListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allApps []calls.SFUListResponse
	for page != nil && len(page.Result) != 0 {
		allApps = append(allApps, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, app := range allApps {
		displayName := app.Name
		if displayName == "" {
			displayName = app.UID
		}
		res := types.Resource{
			Removable:    CallsApp{Client: client.Calls.SFU},
			ResourceID:   app.UID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "CallsApp",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c CallsApp) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, calls.SFUDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
