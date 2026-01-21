package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/stream"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("stream-video", CollectStreamVideos)
}

type StreamVideo struct {
	Client *stream.StreamService
}

func CollectStreamVideos(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.Stream.List(context.TODO(), stream.StreamListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allVideos []stream.Video
	for len(page.Result) != 0 {
		allVideos = append(allVideos, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, video := range allVideos {
		// Use meta.name if available, otherwise use UID
		name := video.UID
		if meta, ok := video.Meta.(map[string]interface{}); ok && meta != nil {
			if metaName, ok := meta["name"]; ok {
				if nameStr, ok := metaName.(string); ok && nameStr != "" {
					name = nameStr
				}
			}
		}

		res := types.Resource{
			Removable:    StreamVideo{Client: client.Stream},
			ResourceID:   video.UID,
			ResourceName: name,
			AccountID:    creds.AccountID,
			ProductName:  "StreamVideo",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c StreamVideo) Remove(accountID string, resourceID string, resourceName string) error {
	return c.Client.Delete(context.TODO(), resourceID, stream.StreamDeleteParams{
		AccountID: cloudflare.F(accountID),
	})
}
