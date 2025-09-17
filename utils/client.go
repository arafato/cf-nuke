package utils

import (
	"github.com/arafato/cf-nuke/types"
	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

func CreateCFClient(creds *types.Credentials) *cloudflare.Client {
	if creds.Mode == "token" {
		return cloudflare.NewClient(
			option.WithAPIToken(creds.APIKey))
	}
	// creds.Mode == "account"
	return cloudflare.NewClient(
		option.WithAPIEmail(creds.User),
		option.WithAPIKey(creds.APIKey))
}
