package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/zero_trust"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("zt-bookmark", CollectZTBookmarks)
}

type ZTBookmark struct {
	Client *zero_trust.AccessBookmarkService
}

func CollectZTBookmarks(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Access.Bookmarks.List(context.TODO(), zero_trust.AccessBookmarkListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		if utils.IsSkippableError(err) {
			utils.AddWarning("ZTBookmark", "", "insufficient permissions or feature not available")
			return nil, nil
		}
		return nil, err
	}

	var allBookmarks []zero_trust.Bookmark
	for page != nil && len(page.Result) != 0 {
		allBookmarks = append(allBookmarks, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, bookmark := range allBookmarks {
		displayName := bookmark.Name
		if displayName == "" {
			displayName = bookmark.Domain
		}
		if displayName == "" {
			displayName = bookmark.ID
		}
		res := types.Resource{
			Removable:    ZTBookmark{Client: client.ZeroTrust.Access.Bookmarks},
			ResourceID:   bookmark.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "ZTBookmark",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c ZTBookmark) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.AccessBookmarkDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
