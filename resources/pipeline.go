package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/pipelines"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("pipeline", CollectPipelines)
}

type Pipeline struct {
	Client *pipelines.PipelineService
}

func CollectPipelines(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	resp, err := client.Pipelines.List(context.TODO(), pipelines.PipelineListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allResources types.Resources
	for _, p := range resp.Results {
		res := types.Resource{
			Removable:    Pipeline{Client: client.Pipelines},
			ResourceID:   p.ID,
			ResourceName: p.Name,
			AccountID:    creds.AccountID,
			ProductName:  "Pipeline",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Pipeline) Remove(accountID string, resourceID string, resourceName string) error {
	// The deprecated API uses pipeline_name for deletion
	return c.Client.Delete(context.TODO(), resourceName, pipelines.PipelineDeleteParams{
		AccountID: cloudflare.F(accountID),
	})
}
