package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v3"
	"github.com/cloudflare/cloudflare-go/v3/option"
	"github.com/cloudflare/cloudflare-go/v3/queues"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
)

func init() {
	infrastructure.RegisterCollector("queue", CollectQueues)
}

type Queue struct {
	Client *cloudflare.Client
}

var ProductName = "Queue"

func CollectQueues(creds *types.Credentials) error {
	client := cloudflare.NewClient(
		option.WithAPIToken(creds.APIKey))

	var allQueues []queues.Queue

	page, err := client.Queues.List(context.TODO(), queues.QueueListParams{
		AccountID: cloudflare.F(creds.AccountId),
	})

	for page != nil {
		allQueues = append(allQueues, page.Result...)
		page, err = page.GetNextPage()
	}
	if err != nil {
		return err
	}

	for _, queue := range allQueues {
		res := types.Resource{
			Removable:   Queue{Client: client},
			ID:          queue.QueueID,
			ProductName: ProductName,
		}

		infrastructure.CollectResource(&res)
	}

	return nil
}

func (Queue) Remove() error {
	return nil
}
