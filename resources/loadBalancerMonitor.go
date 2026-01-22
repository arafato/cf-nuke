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
	infrastructure.RegisterCollector("load-balancer-monitor", CollectLoadBalancerMonitors)
}

type LoadBalancerMonitor struct {
	Client *load_balancers.MonitorService
}

func CollectLoadBalancerMonitors(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	monitorPage, err := client.LoadBalancers.Monitors.List(context.TODO(), load_balancers.MonitorListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})
	if err != nil {
		return nil, err
	}

	var allMonitors []load_balancers.Monitor
	for monitorPage != nil && len(monitorPage.Result) != 0 {
		allMonitors = append(allMonitors, monitorPage.Result...)
		monitorPage, err = monitorPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, monitor := range allMonitors {
		// Use description as name since monitors don't have a Name field
		name := monitor.Description
		if name == "" {
			name = monitor.ID
		}
		res := types.Resource{
			Removable:    LoadBalancerMonitor{Client: client.LoadBalancers.Monitors},
			ResourceID:   monitor.ID,
			ResourceName: name,
			AccountID:    creds.AccountID,
			ProductName:  "LoadBalancerMonitor",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c LoadBalancerMonitor) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, load_balancers.MonitorDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
