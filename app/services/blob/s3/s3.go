package s3

import (
	"bytes"
	"context"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/blob"

	"github.com/getfider/fider/app/pkg/bus"
)

//DefaultClient is an S3 Client
var DefaultClient *s3.S3

func init() {
	bus.Register(&Service{})
}

type Service struct{}

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
		awsSession := session.New(s3Config)
		DefaultClient = s3.New(awsSession)
	}
	bus.AddHandler(s, retrieveBlob)
	bus.AddHandler(s, storeBlob)
	bus.AddHandler(s, deleteBlob)
}

func retrieveBlob(ctx context.Context, cmd *blob.RetrieveBlob) error {
	resp, err := DefaultClient.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(env.Config.BlobStorage.S3.BucketName),
		Key:    aws.String(keyFullPathURL(ctx, cmd.Key)),
	})
	if err != nil {
		if isNotFound(err) {
			return blob.ErrNotFound
		}
		return errors.Wrap(err, "failed to get blob '%s' from S3", cmd.Key)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read blob body '%s' from S3", cmd.Key)
	}

	cmd.Blob = &blob.Blob{
		Content:     bytes,
		ContentType: *resp.ContentType,
		Size:        *resp.ContentLength,
	}
	return nil
}

func storeBlob(ctx context.Context, cmd *blob.StoreBlob) error {
	if err := blob.ValidateKey(cmd.Key); err != nil {
		return errors.Wrap(err, "failed to validate blob key '%s'", cmd.Key)
	}

	reader := bytes.NewReader(cmd.Blob.Content)
	_, err := DefaultClient.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(env.Config.BlobStorage.S3.BucketName),
		Key:         aws.String(keyFullPathURL(ctx, cmd.Key)),
		ContentType: aws.String(cmd.Blob.ContentType),
		ACL:         aws.String(s3.ObjectCannedACLPrivate),
		Body:        reader,
	})
	if err != nil {
		return errors.Wrap(err, "failed to upload blob '%s' to S3", cmd.Key)
	}
	return nil
}

func deleteBlob(ctx context.Context, cmd *blob.DeleteBlob) error {
	_, err := DefaultClient.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(env.Config.BlobStorage.S3.BucketName),
		Key:    aws.String(keyFullPathURL(ctx, cmd.Key)),
	})
	if err != nil && !isNotFound(err) {
		return errors.Wrap(err, "failed to delete blob '%s' from S3", cmd.Key)
	}
	return nil
}

func keyFullPathURL(ctx context.Context, key string) string {
	tenant, ok := ctx.Value(app.TenantCtxKey).(*models.Tenant)
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
