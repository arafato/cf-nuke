package resources

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/custom_certificates"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterZoneCollector("custom-certificate", CollectCustomCertificates)
}

type CustomCertificate struct {
	Client *custom_certificates.CustomCertificateService
	ZoneID string
}

func CollectCustomCertificates(creds *types.Credentials, zone *types.Zone) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	certPage, err := client.CustomCertificates.List(context.TODO(), custom_certificates.CustomCertificateListParams{
		ZoneID: cloudflare.F(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	var allCerts []custom_certificates.CustomCertificate
	for certPage != nil && len(certPage.Result) != 0 {
		allCerts = append(allCerts, certPage.Result...)
		certPage, err = certPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
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

	return allResources, nil
}

func (c CustomCertificate) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, custom_certificates.CustomCertificateDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
