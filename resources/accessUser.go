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
	infrastructure.RegisterAccountCollector("zt-access-user", CollectZTAccessUsers)
}

type ZTAccessUser struct {
	Client  *zero_trust.SeatService
	SeatUID string
}

func CollectZTAccessUsers(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	page, err := client.ZeroTrust.Access.Users.List(context.TODO(), zero_trust.AccessUserListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})
	if err != nil {
		return nil, err
	}

	var allUsers []zero_trust.AccessUserListResponse
	for page != nil && len(page.Result) != 0 {
		allUsers = append(allUsers, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, user := range allUsers {
		// Only include users that occupy a seat
		if !user.AccessSeat && !user.GatewaySeat {
			continue
		}
		displayName := user.Email
		if displayName == "" {
			displayName = user.Name
		}
		if displayName == "" {
			displayName = user.ID
		}
		res := types.Resource{
			Removable:    ZTAccessUser{Client: client.ZeroTrust.Seats, SeatUID: user.SeatUID},
			ResourceID:   user.ID,
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "ZTAccessUser",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c ZTAccessUser) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Edit(context.TODO(), zero_trust.SeatEditParams{
		AccountID: cloudflare.F(accountID),
		Body: []zero_trust.SeatEditParamsBody{
			{
				SeatUID:     cloudflare.F(c.SeatUID),
				AccessSeat:  cloudflare.F(false),
				GatewaySeat: cloudflare.F(false),
			},
		},
	})
	return err
}
