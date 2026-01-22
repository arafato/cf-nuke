package resources

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/workers_for_platforms"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("dispatch-namespace", CollectDispatchNamespaces)
}

type DispatchNamespace struct {
	Client *workers_for_platforms.DispatchNamespaceService
}

func CollectDispatchNamespaces(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.WorkersForPlatforms.Dispatch.Namespaces.List(context.TODO(), workers_for_platforms.DispatchNamespaceListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		// Return empty list for permission/access errors (feature not available for account)
		if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}

	var allNamespaces []workers_for_platforms.DispatchNamespaceListResponse
	for page != nil && len(page.Result) != 0 {
		allNamespaces = append(allNamespaces, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, ns := range allNamespaces {
		res := types.Resource{
			Removable:    DispatchNamespace{Client: client.WorkersForPlatforms.Dispatch.Namespaces},
			ResourceID:   ns.NamespaceID,
			ResourceName: ns.NamespaceName,
			AccountID:    creds.AccountID,
			ProductName:  "DispatchNamespace",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c DispatchNamespace) Remove(accountID string, resourceID string, resourceName string) error {
	// Delete uses namespace name, not ID
	_, err := c.Client.Delete(context.TODO(), resourceName, workers_for_platforms.DispatchNamespaceDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
