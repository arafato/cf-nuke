package resources

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/api_gateway"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterZoneCollector("api-gateway-operation", CollectAPIGatewayOperations)
}

type APIGatewayOperation struct {
	Client *api_gateway.OperationService
	ZoneID string
}

func CollectAPIGatewayOperations(creds *types.Credentials, zone *types.Zone) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	opPage, err := client.APIGateway.Operations.List(context.TODO(), api_gateway.OperationListParams{
		ZoneID: cloudflare.F(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	var allOps []api_gateway.OperationListResponse
	for opPage != nil && len(opPage.Result) != 0 {
		allOps = append(allOps, opPage.Result...)
		opPage, err = opPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, op := range allOps {
		// Create descriptive name: METHOD endpoint
		displayName := fmt.Sprintf("%s %s", op.Method, op.Endpoint)
		res := types.Resource{
			Removable:    APIGatewayOperation{Client: client.APIGateway.Operations, ZoneID: zone.ID},
			ResourceID:   op.OperationID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "APIGatewayOperation",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c APIGatewayOperation) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, api_gateway.OperationDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
