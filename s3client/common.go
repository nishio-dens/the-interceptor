package s3client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"the-interceptor/db"
)

type Client struct {
	s3cli *s3sdk.S3
}

func GetS3Client(bucket *db.S3Bucket) *Client {
	var cred *credentials.Credentials
	if len(bucket.BucketAccessKey) > 0 {
		cred = credentials.NewStaticCredentials(
			bucket.BucketAccessKey,
			bucket.BucketAccessSecret,
			"",
		)
	} else {
		// TODO: Get Auth Info using DetectCred to use IAM
		panic("Not Implemented Yet")
	}

	s3Config := &aws.Config{
		Credentials:      cred,
		Region:           aws.String(bucket.BucketRegion),
		DisableSSL:       aws.Bool(bucket.BucketDisableSsl),
		S3ForcePathStyle: aws.Bool(true),
	}
	if len(bucket.BucketEndpoint) > 0 {
		s3Config.Endpoint = aws.String(bucket.BucketEndpoint)
	}
	newSession := session.New(s3Config)
	cli := s3sdk.New(newSession)

	return &Client{
		s3cli: cli,
	}
}

func (c *Client) ListObjects(input *s3sdk.ListObjectsInput) (*s3sdk.ListObjectsOutput, error) {
	resp, err := c.s3cli.ListObjects(input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) HeadObject(input *s3sdk.HeadObjectInput) (*s3sdk.HeadObjectOutput, error) {
	resp, err := c.s3cli.HeadObject(input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetObject(input *s3sdk.GetObjectInput) (*s3sdk.GetObjectOutput, error) {
	resp, err := c.s3cli.GetObject(input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
