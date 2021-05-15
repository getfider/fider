package s3

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/blob"

	"github.com/getfider/fider/app/pkg/bus"
)

//DefaultClient is an S3 Client
var DefaultClient *s3.S3

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "S3"
}

func (s Service) Category() string {
	return "blobstorage"
}

func (s Service) Enabled() bool {
	return env.Config.BlobStorage.Type == "s3"
}

func (s Service) Init() {
	s3EnvConfig := env.Config.BlobStorage.S3
	if s3EnvConfig.EndpointURL != "" {
		s3Config := &aws.Config{
			Credentials:      credentials.NewStaticCredentials(s3EnvConfig.AccessKeyID, s3EnvConfig.SecretAccessKey, ""),
			Endpoint:         aws.String(s3EnvConfig.EndpointURL),
			Region:           aws.String(s3EnvConfig.Region),
			DisableSSL:       aws.Bool(strings.HasSuffix(s3EnvConfig.EndpointURL, "http://")),
			S3ForcePathStyle: aws.Bool(true),
		}
		awsSession, err := session.NewSession(s3Config)
		if err != nil {
			panic(err)
		}

		DefaultClient = s3.New(awsSession)
	}

	bus.AddHandler(listBlobs)
	bus.AddHandler(getBlobByKey)
	bus.AddHandler(storeBlob)
	bus.AddHandler(deleteBlob)
}

func listBlobs(ctx context.Context, q *query.ListBlobs) error {
	tenant := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	basePath := fmt.Sprintf("tenants/%d/", tenant.ID)
	response, err := DefaultClient.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket:  aws.String(env.Config.BlobStorage.S3.BucketName),
		MaxKeys: aws.Int64(3000),
		Prefix:  aws.String(basePath),
	})
	if err != nil {
		return wrap(err, "failed to list blobs from S3")
	}

	files := make([]string, 0)
	for _, item := range response.Contents {
		key := *item.Key
		files = append(files, key[len(basePath):])
	}

	sort.Strings(files)
	q.Result = files
	return nil
}

func getBlobByKey(ctx context.Context, q *query.GetBlobByKey) error {
	resp, err := DefaultClient.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(env.Config.BlobStorage.S3.BucketName),
		Key:    aws.String(keyFullPathURL(ctx, q.Key)),
	})
	if err != nil {
		if isNotFound(err) {
			return blob.ErrNotFound
		}
		return wrap(err, "failed to get blob '%s' from S3", q.Key)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wrap(err, "failed to read blob body '%s' from S3", q.Key)
	}

	q.Result = &dto.Blob{
		Content:     bytes,
		ContentType: *resp.ContentType,
		Size:        *resp.ContentLength,
	}
	return nil
}

func storeBlob(ctx context.Context, c *cmd.StoreBlob) error {
	if err := blob.ValidateKey(c.Key); err != nil {
		return wrap(err, "failed to validate blob key '%s'", c.Key)
	}

	reader := bytes.NewReader(c.Content)
	_, err := DefaultClient.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(env.Config.BlobStorage.S3.BucketName),
		Key:         aws.String(keyFullPathURL(ctx, c.Key)),
		ContentType: aws.String(c.ContentType),
		ACL:         aws.String(s3.ObjectCannedACLPrivate),
		Body:        reader,
	})
	if err != nil {
		return wrap(err, "failed to upload blob '%s' to S3", c.Key)
	}
	return nil
}

func deleteBlob(ctx context.Context, c *cmd.DeleteBlob) error {
	_, err := DefaultClient.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(env.Config.BlobStorage.S3.BucketName),
		Key:    aws.String(keyFullPathURL(ctx, c.Key)),
	})
	if err != nil && !isNotFound(err) {
		return wrap(err, "failed to delete blob '%s' from S3", c.Key)
	}
	return nil
}

func keyFullPathURL(ctx context.Context, key string) string {
	tenant, ok := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	if ok {
		return path.Join("tenants", strconv.Itoa(tenant.ID), key)
	}
	return key
}

func isNotFound(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		return awsErr.Code() == s3.ErrCodeNoSuchKey
	}
	return false
}

func wrap(err error, format string, a ...interface{}) error {
	if awsErr, ok := err.(awserr.Error); ok {
		return errors.Wrap(awsErr.OrigErr(), format, a...)
	}
	return errors.Wrap(err, format, a...)
}
