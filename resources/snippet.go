package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/snippets"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("snippet", CollectSnippets)
}

type Snippet struct {
	Client *snippets.SnippetService
	ZoneID string
}

func CollectSnippets(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	// First, get all zones for the account
	zonePage, err := client.Zones.List(context.TODO(), zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.F(creds.AccountID)}),
	})
	if err != nil {
		return nil, err
	}

	var allZones []zones.Zone
	for zonePage != nil && len(zonePage.Result) != 0 {
		allZones = append(allZones, zonePage.Result...)
		zonePage, err = zonePage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources

	// For each zone, list all snippets
	for _, zone := range allZones {
		snippetPage, err := client.Snippets.List(context.TODO(), snippets.SnippetListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			if utils.IsSkippableError(err) {
				utils.AddWarning("Snippet", zone.Name, "insufficient permissions or feature not available")
			}
			continue
		}

		var allSnippets []snippets.SnippetListResponse
		for snippetPage != nil && len(snippetPage.Result) != 0 {
			allSnippets = append(allSnippets, snippetPage.Result...)
			snippetPage, err = snippetPage.GetNextPage()
			if err != nil {
				break
			}
		}

		for _, snippet := range allSnippets {
			res := types.Resource{
				Removable:    Snippet{Client: client.Snippets, ZoneID: zone.ID},
				ResourceID:   snippet.SnippetName,
				ResourceName: snippet.SnippetName,
				AccountID:    creds.AccountID,
				ProductName:  "Snippet",
			}
			allResources = append(allResources, &res)
		}
	}

	return allResources, nil
}

func (c Snippet) Remove(accountID string, resourceID string, resourceName string) error {
	// Delete uses snippet name, not ID
	_, err := c.Client.Delete(context.TODO(), resourceID, snippets.SnippetDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
