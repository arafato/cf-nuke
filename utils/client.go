package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/arafato/cf-nuke/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/accounts"
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

func CreateAWSS3Client(accessKeyID, accessKeySecret, accountID string) *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	return client
}

func DeleteTemporaryR2Token(creds *types.Credentials, resources types.Resources) error {
	client := CreateCFClient(creds)
	for _, resource := range resources {
		if resource.ResourceName == types.TemporaryR2TokenName {
			_, err := client.Accounts.Tokens.Delete(context.TODO(), resource.ResourceID, accounts.TokenDeleteParams{
				AccountID: cloudflare.F(creds.AccountID),
			})
			return err
		}
	}
	return nil // usually not reachable since there must always be one hidden resource
}
