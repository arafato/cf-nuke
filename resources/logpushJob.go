package resources

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/logpush"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("logpush-job", CollectLogpushJobs)
}

type LogpushJob struct {
	Client *logpush.JobService
}

func CollectLogpushJobs(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	// Use account-level logpush jobs
	page, err := client.Logpush.Jobs.List(context.TODO(), logpush.JobListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		if utils.IsSkippableError(err) {
			utils.AddWarning("LogpushJob", "", "insufficient permissions or feature not available")
			return nil, nil
		}
		return nil, err
	}

	var allJobs []logpush.LogpushJob
	for page != nil && len(page.Result) != 0 {
		allJobs = append(allJobs, page.Result...)
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	var allResources types.Resources
	for _, job := range allJobs {
		displayName := job.Name
		if displayName == "" {
			displayName = fmt.Sprintf("job-%d", job.ID)
		}
		res := types.Resource{
			Removable:    LogpushJob{Client: client.Logpush.Jobs},
			ResourceID:   fmt.Sprintf("%d", job.ID),
			ResourceName: displayName,
			AccountID:    creds.AccountID,
			ProductName:  "LogpushJob",
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c LogpushJob) Remove(accountID string, resourceID string, resourceName string) error {
	// resourceID is stored as string but API expects int64
	var jobID int64
	fmt.Sscanf(resourceID, "%d", &jobID)

	_, err := c.Client.Delete(context.TODO(), jobID, logpush.JobDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
