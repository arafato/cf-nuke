package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/zero_trust"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterAccountCollector("zt-device-profile", CollectZTDeviceProfiles)
}

type ZTDeviceProfile struct {
	Client *zero_trust.DevicePolicyCustomService
}

func CollectZTDeviceProfiles(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Devices.Policies.Custom.List(context.TODO(), zero_trust.DevicePolicyCustomListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})
	if err != nil {
		return nil, err
	}

	var allResources types.Resources
	for _, profile := range page.Result {
		// Skip the default profile — it cannot be deleted
		if profile.Default {
			continue
		}
		displayName := profile.Name
		if displayName == "" {
			displayName = profile.PolicyID
		}
		res := types.Resource{
			Removable:    ZTDeviceProfile{Client: client.ZeroTrust.Devices.Policies.Custom},
			ResourceID:   profile.PolicyID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "ZTDeviceProfile",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c ZTDeviceProfile) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, zero_trust.DevicePolicyCustomDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
