package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/d1"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("d1", CollectD1)
}

type D1 struct {
	Client *d1.D1Service
}

func CollectD1(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.D1.Database.List(context.TODO(), d1.DatabaseListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	var allD1s []d1.DatabaseListResponse

	if err != nil {
		return nil, err
	}

	for len(page.Result) != 0 {
		allD1s = append(allD1s, page.Result...)
		page, err = page.GetNextPage()
	}

	var allResources types.Resources
	for _, d1 := range allD1s {
		res := types.Resource{
			Removable:    D1{Client: client.D1},
			ResourceID:   d1.UUID,
			ResourceName: d1.Name,
			AccountID:    creds.AccountID,
			ProductName:  "D1",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c D1) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Database.Delete(context.TODO(), resourceID, d1.DatabaseDeleteParams{
		AccountID: cloudflare.F(accountID)})

	return err
}
