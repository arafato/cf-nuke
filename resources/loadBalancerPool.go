package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/load_balancers"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("load-balancer-pool", CollectLoadBalancerPools)
}

type LoadBalancerPool struct {
	Client *load_balancers.PoolService
}

func CollectLoadBalancerPools(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	poolPage, err := client.LoadBalancers.Pools.List(context.TODO(), load_balancers.PoolListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})
	if err != nil {
		return nil, err
	}

	var allPools []load_balancers.Pool
	for len(poolPage.Result) != 0 {
		allPools = append(allPools, poolPage.Result...)
		poolPage, err = poolPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, pool := range allPools {
		res := types.Resource{
			Removable:    LoadBalancerPool{Client: client.LoadBalancers.Pools},
			ResourceID:   pool.ID,
			ResourceName: pool.Name,
			AccountID:    creds.AccountID,
			ProductName:  "LoadBalancerPool",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c LoadBalancerPool) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, load_balancers.PoolDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
