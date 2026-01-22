package resources

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/custom_certificates"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("custom-certificate", CollectCustomCertificates)
}

type CustomCertificate struct {
	Client *custom_certificates.CustomCertificateService
	ZoneID string
}

func CollectCustomCertificates(creds *types.Credentials) (types.Resources, error) {
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

	// For each zone, list all custom certificates
	for _, zone := range allZones {
		certPage, err := client.CustomCertificates.List(context.TODO(), custom_certificates.CustomCertificateListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			if utils.IsSkippableError(err) {
				utils.AddWarning("CustomCertificate", zone.Name, "insufficient permissions")
			}
			continue
		}

		var allCerts []custom_certificates.CustomCertificate
		for certPage != nil && len(certPage.Result) != 0 {
			allCerts = append(allCerts, certPage.Result...)
			certPage, err = certPage.GetNextPage()
			if err != nil {
				break
			}
		}

		for _, cert := range allCerts {
			// Use hosts as display name if available
			displayName := cert.ID
			if len(cert.Hosts) > 0 {
				displayName = strings.Join(cert.Hosts, ", ")
			}
			res := types.Resource{
				Removable:    CustomCertificate{Client: client.CustomCertificates, ZoneID: zone.ID},
				ResourceID:   cert.ID,
				ResourceName: displayName,
				AccountID:    creds.AccountID,
				ProductName:  "CustomCertificate",
			}
			allResources = append(allResources, &res)
		}
	}

	return allResources, nil
}

func (c CustomCertificate) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, custom_certificates.CustomCertificateDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
