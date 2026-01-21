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
	infrastructure.RegisterCollector("stream-live-input", CollectStreamLiveInputs)
}

type StreamLiveInput struct {
	Client *stream.LiveInputService
}

func CollectStreamLiveInputs(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	resp, err := client.Stream.LiveInputs.List(context.TODO(), stream.LiveInputListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allResources types.Resources
	for _, liveInput := range resp.LiveInputs {
		// Use meta.name if available, otherwise use UID
		name := liveInput.UID
		if meta, ok := liveInput.Meta.(map[string]interface{}); ok && meta != nil {
			if metaName, ok := meta["name"]; ok {
				if nameStr, ok := metaName.(string); ok && nameStr != "" {
					name = nameStr
				}
			}
		}

		res := types.Resource{
			Removable:    StreamLiveInput{Client: client.Stream.LiveInputs},
			ResourceID:   liveInput.UID,
			ResourceName: name,
			AccountID:    creds.AccountID,
			ProductName:  "StreamLiveInput",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c StreamLiveInput) Remove(accountID string, resourceID string, resourceName string) error {
	return c.Client.Delete(context.TODO(), resourceID, stream.LiveInputDeleteParams{
		AccountID: cloudflare.F(accountID),
	})
}
