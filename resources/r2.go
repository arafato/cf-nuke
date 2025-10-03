package resources

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/accounts"
	"github.com/cloudflare/cloudflare-go/v6/shared"

	"github.com/arafato/cf-nuke/infrastructure"
	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

func init() {
	infrastructure.RegisterCollector("r2", CollectR2)
}

type R2 struct {
	Client *s3.Client
}

const (
	workersR2StorageWritePermissionGroupId = "bf7481a1826f439697cb59a20b22293e"
)

func CollectR2(creds *types.Credentials) (types.Resources, error) {
	client := utils.CreateCFClient(creds)
	var accessKeyID string
	var tokenValue string

	// Creating temporary account token if we are dealing with deprecated global account creds
	if creds.Mode == types.Account {
		resp, err := client.Accounts.Tokens.New(context.TODO(), accounts.TokenNewParams{
			AccountID: cloudflare.F(creds.AccountID),
			Name:      cloudflare.F(string(types.TemporaryR2TokenName)),
			Policies: cloudflare.F([]shared.TokenPolicyParam{
				{
					Effect: cloudflare.F(shared.TokenPolicyEffectAllow),
					PermissionGroups: cloudflare.F([]shared.TokenPolicyPermissionGroupParam{
						{
							ID: cloudflare.F(workersR2StorageWritePermissionGroupId),
						},
					}),
					Resources: cloudflare.F[shared.TokenPolicyResourcesUnionParam](shared.TokenPolicyResourcesIAMResourcesTypeObjectStringParam(map[string]string{
						fmt.Sprintf("com.cloudflare.api.account.%s", creds.AccountID): "*",
					})),
				}}),
		})

		if err != nil {
			return nil, err
		}

		accessKeyID = resp.ID
		tokenValue = resp.Value
	}

	if creds.Mode == types.Token {
		resp, err := client.User.Tokens.Verify(context.TODO()) // assume sufficient permissions
		if err != nil {
			return nil, err
		}

		accessKeyID = resp.ID
		tokenValue = creds.APIKey
	}

	hash := sha256.Sum256([]byte(tokenValue))
	accessKeySecret := hex.EncodeToString(hash[:])

	time.Sleep(3 * time.Second) //if we don't wait the token is not yet ready for S3 authentication
	s3Client := utils.CreateAWSS3Client(accessKeyID, accessKeySecret, creds.AccountID)

	output, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	for _, object := range output.Buckets {
		obj, _ := json.MarshalIndent(object, "", "\t")
		fmt.Println(string(obj))
	}

	var allResources types.Resources
	for _, b := range output.Buckets {
		res := types.Resource{
			Removable:    R2{Client: s3Client},
			ResourceID:   *b.Name,
			ResourceName: *b.Name,
			AccountID:    creds.AccountID,
			ProductName:  "R2",
			State:        types.Ready,
		}
		allResources = append(allResources, &res)
	}

	return allResources, nil
}

func (c R2) Remove(accountID string, resourceID string, resourceName string) error {

	listResp, err := c.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(resourceID),
	})
	if err != nil {
		return err
	}

	// 2. Delete objects (in batches)
	if len(listResp.Contents) > 0 {
		objects := make([]s3types.ObjectIdentifier, len(listResp.Contents))
		for i, obj := range listResp.Contents {
			objects[i] = s3types.ObjectIdentifier{Key: obj.Key}
		}

		_, err = c.Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
			Bucket: aws.String(resourceID),
			Delete: &s3types.Delete{
				Objects: objects,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return err
		}
	}

	// 3. Delete the bucket
	_, err = c.Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(resourceID),
	})
	if err != nil {
		return err
	}

	// TODO: Delete temporary token
	return nil
}
