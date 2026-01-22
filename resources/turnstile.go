package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/turnstile"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("turnstile", CollectTurnstile)
}

type Turnstile struct {
	Client *turnstile.WidgetService
}

func CollectTurnstile(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Turnstile.Widgets.List(context.TODO(), turnstile.WidgetListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allWidgets []turnstile.WidgetListResponse
	for page != nil && len(page.Result) != 0 {
		allWidgets = append(allWidgets, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, widget := range allWidgets {
		res := types.Resource{
			Removable:    Turnstile{Client: client.Turnstile.Widgets},
			ResourceID:   widget.Sitekey,
			ResourceName: widget.Name,
			AccountID:    creds.AccountID,
			ProductName:  "Turnstile",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Turnstile) Remove(accountID string, resourceID string, resourceName string) error {
	// Delete uses sitekey as identifier
	_, err := c.Client.Delete(context.TODO(), resourceID, turnstile.WidgetDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
