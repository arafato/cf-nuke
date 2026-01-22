package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/waiting_rooms"
	"github.com/cloudflare/cloudflare-go/v6/zones"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("waiting-room", CollectWaitingRooms)
}

type WaitingRoom struct {
	Client *waiting_rooms.WaitingRoomService
	ZoneID string
}

func CollectWaitingRooms(creds *types.Credentials) (types.Resources, error) {
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

	// For each zone, list all waiting rooms
	for _, zone := range allZones {
		wrPage, err := client.WaitingRooms.List(context.TODO(), waiting_rooms.WaitingRoomListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			// Skip zones where we might not have permissions
			continue
		}

		var allWaitingRooms []waiting_rooms.WaitingRoom
		for wrPage != nil && len(wrPage.Result) != 0 {
			allWaitingRooms = append(allWaitingRooms, wrPage.Result...)
			wrPage, err = wrPage.GetNextPage()
			if err != nil {
				break
			}
		}

		for _, wr := range allWaitingRooms {
			res := types.Resource{
				Removable:    WaitingRoom{Client: client.WaitingRooms, ZoneID: zone.ID},
				ResourceID:   wr.ID,
				ResourceName: wr.Name,
				AccountID:    creds.AccountID,
				ProductName:  "WaitingRoom",
				State:        types.Ready,
			}
			allResources = append(allResources, &res)
		}
	}

	return allResources, nil
}

func (c WaitingRoom) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, waiting_rooms.WaitingRoomDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
