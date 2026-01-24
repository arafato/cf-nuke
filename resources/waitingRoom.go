package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/waiting_rooms"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterZoneCollector("waiting-room", CollectWaitingRooms)
}

type WaitingRoom struct {
	Client *waiting_rooms.WaitingRoomService
	ZoneID string
}

func CollectWaitingRooms(creds *types.Credentials, zone *types.Zone) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	wrPage, err := client.WaitingRooms.List(context.TODO(), waiting_rooms.WaitingRoomListParams{
		ZoneID: cloudflare.F(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	var allWaitingRooms []waiting_rooms.WaitingRoom
	for wrPage != nil && len(wrPage.Result) != 0 {
		allWaitingRooms = append(allWaitingRooms, wrPage.Result...)
		wrPage, err = wrPage.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, wr := range allWaitingRooms {
		res := types.Resource{
			Removable:    WaitingRoom{Client: client.WaitingRooms, ZoneID: zone.ID},
			ResourceID:   wr.ID,
			ResourceName: wr.Name,
			AccountID:    creds.AccountID,
			ProductName:  "WaitingRoom",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c WaitingRoom) Remove(accountID string, resourceID string, resourceName string) error {
	_, err := c.Client.Delete(context.TODO(), resourceID, waiting_rooms.WaitingRoomDeleteParams{
		ZoneID: cloudflare.F(c.ZoneID),
	})

	return err
}
