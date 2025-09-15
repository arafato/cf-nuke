package resources

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v3"
	"github.com/cloudflare/cloudflare-go/v3/option"
	"github.com/cloudflare/cloudflare-go/v3/queues"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
)

func init() {
	infrastructure.RegisterCollector("queue", CollectQueues)
}

type Queue struct{}

func CollectQueues(creds *types.Credentials) ([]types.Resource, error) {
	client := cloudflare.NewClient(
		option.WithAPIToken(creds.APIKey))

	page, err := client.Queues.List(context.TODO(), queues.QueueListParams{
		AccountID: cloudflare.F(creds.AccountId),
	})

	var res = types.Resource{
		Removable:   Queue{},
		ID:          "id",
		ProductName: "Queue",
	}

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", page)
	// create a
	return nil, nil
}

func (Queue) Remove() error {
	return nil
}
