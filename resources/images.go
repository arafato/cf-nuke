package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/images"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterAccountCollector("images", CollectImages)
}

type Image struct {
	Client *images.V1Service
}

func CollectImages(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Images.V1.List(context.TODO(), images.V1ListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})
	if err != nil {
		return nil, err
	}

	var allImages []images.Image
	for page != nil && len(page.Result.Items) != 0 {
		for _, item := range page.Result.Items {
			allImages = append(allImages, item.Images...)
		}
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, img := range allImages {
		displayName := img.Filename
		if displayName == "" {
			displayName = img.ID
		}
		res := types.Resource{
			Removable:    Image{Client: client.Images.V1},
			ResourceID:   img.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "Image",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c Image) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, images.V1DeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
