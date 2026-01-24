package resources

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/mtls_certificates"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterAccountCollector("mtls-certificate", CollectMTLSCertificates)
}

type MTLSCertificate struct {
	Client *mtls_certificates.MTLSCertificateService
}

func CollectMTLSCertificates(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.MTLSCertificates.List(context.TODO(), mtls_certificates.MTLSCertificateListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	var allCerts []mtls_certificates.MTLSCertificate
	for page != nil && len(page.Result) != 0 {
		allCerts = append(allCerts, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, cert := range allCerts {
		// Skip Cloudflare-managed certificates (e.g., "Gateway CA - Cloudflare Managed G1")
		if strings.Contains(cert.Name, "Cloudflare Managed") {
			continue
		}

		displayName := cert.Name
		if displayName == "" {
			displayName = cert.ID
		}
		res := types.Resource{
			Removable:    MTLSCertificate{Client: client.MTLSCertificates},
			ResourceID:   cert.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "MTLSCertificate",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c MTLSCertificate) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, mtls_certificates.MTLSCertificateDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
