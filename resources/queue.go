package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/queues"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
)

func init() {
	infrastructure.RegisterCollector("queue", CollectQueues)
}

type Queue struct {
	Client *queues.QueueService
}

func CollectQueues(creds *types.Credentials) (types.Resources, error) {
	client := cloudflare.NewClient(
		option.WithAPIToken(creds.APIKey))

	var allQueues []queues.Queue

	page, err := client.Queues.List(context.TODO(), queues.QueueListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	for page != nil {
		allQueues = append(allQueues, page.Result...)
		page, err = page.GetNextPage()
	}

	var allResources types.Resources
	for _, queue := range allQueues {
		res := types.Resource{
			Removable:    Queue{Client: client.Queues},
			ResourceID:   queue.QueueID,
			ResourceName: queue.QueueName,
			AccountID:    creds.AccountID,
			ProductName:  "Queue",
			State:        types.Ready,
		}

		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (q Queue) Remove(accountID string, resourceID string) error {
	_, err := q.Client.Delete(context.TODO(), resourceID, queues.QueueDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
