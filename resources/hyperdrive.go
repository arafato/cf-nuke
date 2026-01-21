package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/hyperdrive"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("hyperdrive", CollectHyperdrive)
}

type Hyperdrive struct {
	Client *hyperdrive.ConfigService
}

func CollectHyperdrive(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Hyperdrive.Configs.List(context.TODO(), hyperdrive.ConfigListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allConfigs []hyperdrive.Hyperdrive
	for len(page.Result) != 0 {
		allConfigs = append(allConfigs, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, cfg := range allConfigs {
		res := types.Resource{
			Removable:    Hyperdrive{Client: client.Hyperdrive.Configs},
			ResourceID:   cfg.ID,
			ResourceName: cfg.Name,
			AccountID:    creds.AccountID,
			ProductName:  "Hyperdrive",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Hyperdrive) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, hyperdrive.ConfigDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
