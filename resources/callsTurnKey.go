package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/calls"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterAccountCollector("calls-turn-key", CollectCallsTurnKeys)
}

type CallsTurnKey struct {
	Client *calls.TURNService
}

func CollectCallsTurnKeys(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Calls.TURN.List(context.TODO(), calls.TURNListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allKeys []calls.TURNListResponse
	for page != nil && len(page.Result) != 0 {
		allKeys = append(allKeys, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, key := range allKeys {
		displayName := key.Name
		if displayName == "" {
			displayName = key.UID
		}
		res := types.Resource{
			Removable:    CallsTurnKey{Client: client.Calls.TURN},
			ResourceID:   key.UID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "CallsTurnKey",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c CallsTurnKey) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, calls.TURNDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
