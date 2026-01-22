package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/keyless_certificates"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("keyless-certificate", CollectKeylessCertificates)
}

type KeylessCertificate struct {
	Client *keyless_certificates.KeylessCertificateService
	ZoneID string
}

func CollectKeylessCertificates(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	// First, get all zones for the account
	zonePage, err := client.Zones.List(context.TODO(), zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.F(creds.AccountID)}),
	})
	if err != nil {
		return nil, err
	}

	var allZones []zones.Zone
	for zonePage != nil && len(zonePage.Result) != 0 {
		allZones = append(allZones, zonePage.Result...)
		zonePage, err = zonePage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources

	// For each zone, list all keyless certificates
	for _, zone := range allZones {
		certPage, err := client.KeylessCertificates.List(context.TODO(), keyless_certificates.KeylessCertificateListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			if utils.IsSkippableError(err) {
				utils.AddWarning("KeylessCertificate", zone.Name, "insufficient permissions")
			}
			continue
		}

		var allCerts []keyless_certificates.KeylessCertificate
		for certPage != nil && len(certPage.Result) != 0 {
			allCerts = append(allCerts, certPage.Result...)
			certPage, err = certPage.GetNextPage()
			if err != nil {
				break
			}
		}

		for _, cert := range allCerts {
			displayName := cert.Name
			if displayName == "" {
				displayName = cert.Host
			}
			if displayName == "" {
				displayName = cert.ID
			}
			res := types.Resource{
				Removable:    KeylessCertificate{Client: client.KeylessCertificates, ZoneID: zone.ID},
				ResourceID:   cert.ID,
				ResourceName: displayName,
				AccountID:    creds.AccountID,
				ProductName:  "KeylessCertificate",
			}
			allResources = append(allResources, &res)
		}
	}

	return allResources, nil
}

func (c KeylessCertificate) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, keyless_certificates.KeylessCertificateDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
