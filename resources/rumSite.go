package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/rum"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("rum-site", CollectRUMSites)
}

type RUMSite struct {
	Client *rum.SiteInfoService
}

func CollectRUMSites(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.RUM.SiteInfo.List(context.TODO(), rum.SiteInfoListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		if utils.IsSkippableError(err) {
			utils.AddWarning("RUMSite", "", "insufficient permissions or feature not available")
			return nil, nil
		}
		return nil, err
	}

	var allSites []rum.Site
	for page != nil && len(page.Result) != 0 {
		allSites = append(allSites, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, site := range allSites {
		res := types.Resource{
			Removable:    RUMSite{Client: client.RUM.SiteInfo},
			ResourceID:   site.SiteTag,
			ResourceName: site.SiteTag,
			AccountID:    creds.AccountID,
			ProductName:  "RUMSite",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c RUMSite) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, rum.SiteInfoDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
