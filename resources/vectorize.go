package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/vectorize"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("vectorize", CollectVectorize)
}

type Vectorize struct {
	Client *vectorize.IndexService
}

func CollectVectorize(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Vectorize.Indexes.List(context.TODO(), vectorize.IndexListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allIndexes []vectorize.CreateIndex
	for page != nil && len(page.Result) != 0 {
		allIndexes = append(allIndexes, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, idx := range allIndexes {
		res := types.Resource{
			Removable:    Vectorize{Client: client.Vectorize.Indexes},
			ResourceID:   idx.Name,
			ResourceName: idx.Name,
			AccountID:    creds.AccountID,
			ProductName:  "Vectorize",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Vectorize) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceName, vectorize.IndexDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
