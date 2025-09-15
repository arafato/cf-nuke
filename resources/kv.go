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

func CollectKV(creds *types.Credentials) error {
	client := cloudflare.NewClient(
		option.WithAPIToken(creds.APIKey),
	)

	page, err := client.KV.Namespaces.List(context.TODO(), kv.NamespaceListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	var allKVs []kv.Namespace

	// page is never nil
	for len(page.Result) != 0 {
		allKVs = append(allKVs, page.Result...)
		page, err = page.GetNextPage()
	}
	if err != nil {
		return err
	}

	for _, kv := range allKVs {
		res := types.Resource{
			Removable:    KV{Client: client.KV},
			ResourceID:   kv.ID,
			ResourceName: kv.Title,
			AccountID:    creds.AccountID,
			ProductName:  "KV",
		}
		infrastructure.CollectResource(&res)
	}

	return nil
}

func (q KV) Remove(accountID string, resourceID string) error {
	return nil
}
