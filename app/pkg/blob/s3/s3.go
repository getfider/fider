package s3

import (
	"bytes"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

var _ blob.Storage = (*Storage)(nil)

// Storage stores blobs on an S3 compatible service
type Storage struct {
	bucket *string
	tenant *models.Tenant
}

func isNotFound(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		return awsErr.Code() == s3.ErrCodeNoSuchKey
	}
	return false
}

//DefaultClient is an S3 Client
var DefaultClient *s3.S3

func init() {
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
}

// NewStorage creates a S3 compatible service storage
func NewStorage(bucket string) *Storage {
	return &Storage{
		bucket: aws.String(bucket),
	}
}

func (s *Storage) keyFullPathURL(key string) string {
	if s.tenant != nil {
		return path.Join("tenants", strconv.Itoa(s.tenant.ID), key)
	}
	return key
}

// SetCurrentTenant to current context
func (s *Storage) SetCurrentTenant(tenant *models.Tenant) {
	s.tenant = tenant
}

// Get returns a blob with given key
func (s *Storage) Get(key string) (*blob.Blob, error) {
	resp, err := DefaultClient.GetObject(&s3.GetObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(s.keyFullPathURL(key)),
	})
	if err != nil {
		if isNotFound(err) {
			return nil, blob.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to get blob '%s' from S3", key)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read blob body '%s' from S3", key)
	}

	return &blob.Blob{
		Key:         key,
		Size:        *resp.ContentLength,
		ContentType: *resp.ContentType,
		Object:      bytes,
	}, nil
}

// Delete a blob with given key
func (s *Storage) Delete(key string) error {
	_, err := DefaultClient.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(s.keyFullPathURL(key)),
	})
	if err != nil && !isNotFound(err) {
		return errors.Wrap(err, "failed to delete blob '%s' from S3", key)
	}
	return nil
}

// Put a blob with given key and content. Blobs with same key are replaced.
func (s *Storage) Put(key string, content []byte, contentType string) error {
	if err := blob.ValidateKey(key); err != nil {
		return errors.Wrap(err, "failed to validate blob key '%s'", key)
	}

	reader := bytes.NewReader(content)
	_, err := DefaultClient.PutObject(&s3.PutObjectInput{
		Bucket:      s.bucket,
		Key:         aws.String(s.keyFullPathURL(key)),
		ContentType: aws.String(contentType),
		ACL:         aws.String(s3.ObjectCannedACLPrivate),
		Body:        reader,
	})
	if err != nil {
		return errors.Wrap(err, "failed to upload blob '%s' to S3", key)
	}
	return nil
}
