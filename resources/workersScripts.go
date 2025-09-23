package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/workers"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("workers-scripts", CollectWorkers)
}

type WorkersScripts struct {
	Client *workers.ScriptService
}

func CollectWorkers(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Workers.Scripts.List(context.TODO(), workers.ScriptListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	var allWorkersScripts []workers.Script

	if err != nil {
		return nil, err
	}

	for page != nil {
		allWorkersScripts = append(allWorkersScripts, page.Result...)
		page, err = page.GetNextPage()
	}

	var allResources types.Resources
	for _, script := range allWorkersScripts {
		res := types.Resource{
			Removable:    WorkersScripts{Client: client.Workers.Scripts},
			ResourceID:   script.ID,
			ResourceName: script.ID,
			AccountID:    creds.AccountID,
			ProductName:  "workers-scripts",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c WorkersScripts) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, workers.ScriptDeleteParams{
		AccountID: cloudflare.F(accountID)})

	return err
}
