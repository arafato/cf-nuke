package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/r2"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
)

func init() {
	infrastructure.RegisterCollector("r2", CollectR2)
}

type R2 struct {
	Client *r2.R2Service
}

func CollectR2(creds *types.Credentials) (types.Resources, error) {
	client := cloudflare.NewClient(
		option.WithAPIToken(creds.APIKey),
	)

	resp, err := client.R2.Buckets.List(context.TODO(), r2.BucketListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allResources types.Resources
	for _, bucket := range resp.Buckets {
		res := types.Resource{
			Removable:    R2{Client: client.R2},
			ResourceID:   bucket.Name,
			ResourceName: bucket.Name,
			AccountID:    creds.AccountID,
			ProductName:  "R2",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c R2) Remove(accountID string, resourceID string) error {
	_, err := c.Client.Buckets.Delete(context.TODO(), resourceID, r2.BucketDeleteParams{
		AccountID: cloudflare.F(accountID)})

	return err
}
