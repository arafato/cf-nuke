package resources

import (
	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
)

func init() {
	infrastructure.Register("queue", ListQueues)
}

func ListQueues() ([]types.Resource, error) {
	return nil, nil
}
