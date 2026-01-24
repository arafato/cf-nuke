package resources

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterZoneCollector("dns-record", CollectDNSRecords)
}

type DNSRecord struct {
	Client *dns.RecordService
	ZoneID string
}

func CollectDNSRecords(creds *types.Credentials, zone *types.Zone) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	recordPage, err := client.DNS.Records.List(context.TODO(), dns.RecordListParams{
		ZoneID: cloudflare.F(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	var allRecords []dns.RecordResponse
	for recordPage != nil && len(recordPage.Result) != 0 {
		allRecords = append(allRecords, recordPage.Result...)
		recordPage, err = recordPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, record := range allRecords {
		// Create a descriptive name: TYPE name (e.g., "A example.com")
		displayName := fmt.Sprintf("%s %s", record.Type, record.Name)
		res := types.Resource{
			Removable:    DNSRecord{Client: client.DNS.Records, ZoneID: zone.ID},
			ResourceID:   record.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "DNSRecord",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c DNSRecord) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, dns.RecordDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
