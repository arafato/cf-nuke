package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/workflows"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("workflow", CollectWorkflows)
}

type Workflow struct {
	Client *workflows.WorkflowService
}

func CollectWorkflows(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Workflows.List(context.TODO(), workflows.WorkflowListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	var allWorkflows []workflows.WorkflowListResponse

	if err != nil {
		return nil, err
	}

	for page != nil && len(page.Result) != 0 {
		allWorkflows = append(allWorkflows, page.Result...)
		page, err = page.GetNextPage()
	}

	var allResources types.Resources
	for _, wf := range allWorkflows {
		res := types.Resource{
			Removable:    Workflow{Client: client.Workflows},
			ResourceID:   wf.ID,
			ResourceName: wf.Name,
			AccountID:    creds.AccountID,
			ProductName:  "Workflow",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Workflow) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceName, workflows.WorkflowDeleteParams{
		AccountID: cloudflare.F(accountID)})

	return err
}
