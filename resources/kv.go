package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v3"
	"github.com/cloudflare/cloudflare-go/v3/kv"
	"github.com/cloudflare/cloudflare-go/v3/option"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
)

func init() {
	infrastructure.RegisterCollector("kv", CollectKV)
}

type KV struct {
	Client *kv.KVService
}

func CollectKV(creds *types.Credentials) (types.Resources, error) {
	client := cloudflare.NewClient(
		option.WithAPIToken(creds.APIKey),
	)

	page, err := client.KV.Namespaces.List(context.TODO(), kv.NamespaceListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	var allKVs []kv.Namespace

	if err != nil {
		return nil, err
	}

	for len(page.Result) != 0 {
		allKVs = append(allKVs, page.Result...)
		page, err = page.GetNextPage()
	}

	var allResources types.Resources
	for _, kv := range allKVs {
		res := types.Resource{
			Removable:    KV{Client: client.KV},
			ResourceID:   kv.ID,
			ResourceName: kv.Title,
			AccountID:    creds.AccountID,
			ProductName:  "KV",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (q KV) Remove(accountID string, resourceID string) error {
	return nil
}
