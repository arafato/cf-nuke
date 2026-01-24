package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/secrets_store"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterAccountCollector("secrets-store", CollectSecretsStore)
}

type SecretsStore struct {
	Client *secrets_store.SecretsStoreService
}

func CollectSecretsStore(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.SecretsStore.Stores.List(context.TODO(), secrets_store.StoreListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	var allSecretStores []secrets_store.StoreListResponse

	if err != nil {
		return nil, err
	}

	for len(page.Result) != 0 {
		allSecretStores = append(allSecretStores, page.Result...)
		page, err = page.GetNextPage()
	}

	var allResources types.Resources
	for _, store := range allSecretStores {
		res := types.Resource{
			Removable:    SecretsStore{Client: client.SecretsStore},
			ResourceID:   store.ID,
			ResourceName: store.Name,
			AccountID:    creds.AccountID,
			ProductName:  "SecretsStore",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c SecretsStore) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Stores.Delete(context.TODO(), resourceID, secrets_store.StoreDeleteParams{
		AccountID: cloudflare.F(accountID)})

	return err
}
