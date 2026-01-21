package resources

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/pages"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("pages-project", CollectPagesProjects)
}

type PagesProject struct {
	Client *pages.ProjectService
}

func CollectPagesProjects(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)

	// The SDK List returns Deployment objects, but they contain ProjectID and ProjectName
	// We need to deduplicate by project name since multiple deployments can belong to the same project
	page, err := client.Pages.Projects.List(context.TODO(), pages.ProjectListParams{
		AccountID: cloudflare.F(creds.AccountID),
	})

	if err != nil {
		return nil, err
	}

	// Use a map to deduplicate projects by name
	projectsSeen := make(map[string]bool)
	var allResources types.Resources

	for len(page.Result) != 0 {
		for _, deployment := range page.Result {
			// Skip if we've already seen this project
			if projectsSeen[deployment.ProjectName] {
				continue
			}
			projectsSeen[deployment.ProjectName] = true

			res := types.Resource{
				Removable:    PagesProject{Client: client.Pages.Projects},
				ResourceID:   deployment.ProjectID,
				ResourceName: deployment.ProjectName,
				AccountID:    creds.AccountID,
				ProductName:  "PagesProject",
				State:        types.Ready,
			}
			allResources = append(allResources, &res)
		}
		page, err = page.GetNextPage()
		if err != nil {
			break
		}
	}

	return allResources, nil
}

func (c PagesProject) Remove(accountID string, resourceID string, resourceName string) error {
	// Delete uses project name, not ID
	_, err := c.Client.Delete(context.TODO(), resourceName, pages.ProjectDeleteParams{
		AccountID: cloudflare.F(accountID),
	})

	return err
}
