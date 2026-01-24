package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/client_certificates"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterZoneCollector("client-certificate", CollectClientCertificates)
}

type ClientCertificate struct {
	Client *client_certificates.ClientCertificateService
	ZoneID string
}

func CollectClientCertificates(creds *types.Credentials, zone *types.Zone) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	certPage, err := client.ClientCertificates.List(context.TODO(), client_certificates.ClientCertificateListParams{
		ZoneID: cloudflare.F(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	var allCerts []client_certificates.ClientCertificate
	for certPage != nil && len(certPage.Result) != 0 {
		allCerts = append(allCerts, certPage.Result...)
		certPage, err = certPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, cert := range allCerts {
		displayName := cert.CommonName
		if displayName == "" {
			displayName = cert.ID
		}
		res := types.Resource{
			Removable:    ClientCertificate{Client: client.ClientCertificates, ZoneID: zone.ID},
			ResourceID:   cert.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "ClientCertificate",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c ClientCertificate) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, client_certificates.ClientCertificateDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
