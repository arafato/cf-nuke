package resources

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/origin_ca_certificates"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterZoneCollector("origin-ca-certificate", CollectOriginCACertificates)
}

type OriginCACertificate struct {
	Client *origin_ca_certificates.OriginCACertificateService
}

func CollectOriginCACertificates(creds *types.Credentials, zone *types.Zone) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	certPage, err := client.OriginCACertificates.List(context.TODO(), origin_ca_certificates.OriginCACertificateListParams{
		ZoneID: cloudflare.F(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	var allCerts []origin_ca_certificates.OriginCACertificate
	for certPage != nil && len(certPage.Result) != 0 {
		allCerts = append(allCerts, certPage.Result...)
		certPage, err = certPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, cert := range allCerts {
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

	return allResources, nil
}

func (c OriginCACertificate) Remove(accountID string, resourceID string, resourceName string) error {
	// Origin CA certificate delete doesn't require zone ID - it's a global operation
	_, err := c.Client.Delete(context.TODO(), resourceID)

	return err
}
