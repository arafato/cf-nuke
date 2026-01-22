package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"

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

func CollectR2(creds *types.Credentials) (types.Resources, error) {
	s3Client := utils.CreateAWSS3Client(creds.S3AccessKeyID, creds.S3AccessSecret, creds.AccountID)

	output, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var allResources types.Resources
	for _, b := range output.Buckets {
		res := types.Resource{
			Removable:    R2{Client: s3Client},
			ResourceID:   *b.Name,
			ResourceName: *b.Name,
			AccountID:    creds.AccountID,
			ProductName:  "R2",
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

	return nil
}
