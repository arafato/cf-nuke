package resources

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/origin_ca_certificates"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("origin-ca-certificate", CollectOriginCACertificates)
}

type OriginCACertificate struct {
	Client *origin_ca_certificates.OriginCACertificateService
}

func CollectOriginCACertificates(creds *types.Credentials) (types.Resources, error) {
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

	// Track seen certificate IDs to avoid duplicates (same cert can be in multiple zones)
	seenCerts := make(map[string]bool)
	var allResources types.Resources

	// For each zone, list all origin CA certificates
	for _, zone := range allZones {
		certPage, err := client.OriginCACertificates.List(context.TODO(), origin_ca_certificates.OriginCACertificateListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			if utils.IsSkippableError(err) {
				utils.AddWarning("OriginCACertificate", zone.Name, "insufficient permissions")
			}
			continue
		}

		var allCerts []origin_ca_certificates.OriginCACertificate
		for certPage != nil && len(certPage.Result) != 0 {
			allCerts = append(allCerts, certPage.Result...)
			certPage, err = certPage.GetNextPage()
			if err != nil {
				break
			}
		}

		for _, cert := range allCerts {
			// Skip duplicates
			if seenCerts[cert.ID] {
				continue
			}
			seenCerts[cert.ID] = true

			// Use hostnames as display name if available
			displayName := cert.ID
			if len(cert.Hostnames) > 0 {
				displayName = strings.Join(cert.Hostnames, ", ")
			}
			res := types.Resource{
				Removable:    OriginCACertificate{Client: client.OriginCACertificates},
				ResourceID:   cert.ID,
				ResourceName: displayName,
				AccountID:    creds.AccountID,
				ProductName:  "OriginCACertificate",
			}
			allResources = append(allResources, &res)
		}
	}

	return allResources, nil
}

func (c OriginCACertificate) Remove(accountID string, resourceID string, resourceName string) error {
	// Origin CA certificate delete doesn't require zone ID - it's a global operation
	_, err := c.Client.Delete(context.TODO(), resourceID)

	return err
}
