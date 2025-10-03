package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	"github.com/cloudflare/cloudflare-go/v6/shared"
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

func CreateTemporaryR2Token(creds *types.Credentials) (string, string, error) {
	client := CreateCFClient(creds)
	var accessKeyID string
	var tokenValue string

	if creds.Mode == types.Account {

		resp, err := client.Accounts.Tokens.New(context.TODO(), accounts.TokenNewParams{
			AccountID: cloudflare.F(creds.AccountID),
			Name:      cloudflare.F(string(types.TemporaryR2TokenName)),
			Policies: cloudflare.F([]shared.TokenPolicyParam{
				{
					Effect: cloudflare.F(shared.TokenPolicyEffectAllow),
					PermissionGroups: cloudflare.F([]shared.TokenPolicyPermissionGroupParam{
						{
							ID: cloudflare.F(types.WorkersR2StorageWritePermissionGroupId),
						},
					}),
					Resources: cloudflare.F[shared.TokenPolicyResourcesUnionParam](shared.TokenPolicyResourcesIAMResourcesTypeObjectStringParam(map[string]string{
						fmt.Sprintf("com.cloudflare.api.account.%s", creds.AccountID): "*",
					})),
				}}),
		})

		if err != nil {
			return "", "", err
		}

		accessKeyID = resp.ID
		tokenValue = resp.Value
	}

	if creds.Mode == types.Token {
		resp, err := client.User.Tokens.Verify(context.TODO()) // assume sufficient permissions
		if err != nil {
			return "", "", err
		}

		accessKeyID = resp.ID
		tokenValue = creds.APIKey
	}

	hash := sha256.Sum256([]byte(tokenValue))
	accessKeySecret := hex.EncodeToString(hash[:])
	return accessKeyID, accessKeySecret, nil
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
	return nil
}
