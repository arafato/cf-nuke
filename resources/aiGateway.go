package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/ai_gateway"
	"github.com/cloudflare/cloudflare-go/v6/option"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
)

func init() {
	infrastructure.RegisterCollector("ai-gateway", CollectAIGateway)
}

type AIGateway struct {
	Client *ai_gateway.AIGatewayService
}

func CollectAIGateway(creds *types.Credentials) (types.Resources, error) {
	client := cloudflare.NewClient(
		option.WithAPIToken(creds.APIKey),
	)

	page, err := client.AIGateway.List(context.TODO(), ai_gateway.AIGatewayListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	var allAIGateways []ai_gateway.AIGatewayListResponse

	if err != nil {
		return nil, err
	}

	for len(page.Result) != 0 {
		allAIGateways = append(allAIGateways, page.Result...)
		page, err = page.GetNextPage()
	}

	var allResources types.Resources
	for _, aigw := range allAIGateways {
		res := types.Resource{
			Removable:    AIGateway{Client: client.AIGateway},
			ResourceID:   aigw.ID,
			ResourceName: aigw.ID,
			AccountID:    creds.AccountID,
			ProductName:  "AIGateway",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c AIGateway) Remove(accountID string, resourceID string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, ai_gateway.AIGatewayDeleteParams{
		AccountID: cloudflare.F(accountID)})

	return err
}
