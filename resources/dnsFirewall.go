package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns_firewall"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterAccountCollector("dns-firewall", CollectDNSFirewalls)
}

type DNSFirewall struct {
	Client *dns_firewall.DNSFirewallService
}

func CollectDNSFirewalls(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.DNSFirewall.List(context.TODO(), dns_firewall.DNSFirewallListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allClusters []dns_firewall.DNSFirewallListResponse
	for page != nil && len(page.Result) != 0 {
		allClusters = append(allClusters, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, cluster := range allClusters {
		displayName := cluster.Name
		if displayName == "" {
			displayName = cluster.ID
		}
		res := types.Resource{
			Removable:    DNSFirewall{Client: client.DNSFirewall},
			ResourceID:   cluster.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "DNSFirewall",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c DNSFirewall) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, dns_firewall.DNSFirewallDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
