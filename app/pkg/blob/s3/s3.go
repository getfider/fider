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
	"github.com/getfider/fider/app/pkg/errors"
)

var _ blob.Storage = (*Storage)(nil)

// Storage stores blobs on an S3 compatible service
type Storage struct {
	awsSession *session.Session
	bucket     *string
}

// Session is a per-request object to interact with the storage
type Session struct {
	storage  *Storage
	tenant   *models.Tenant
	s3Client *s3.S3
}

func isNotFound(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		return awsErr.Code() == s3.ErrCodeNoSuchKey
	}
	return false
}

// NewStorage creates a S3 compatible service storage
func NewStorage(endpointURL, region, accessKeyID, secretAccessKey, bucket string) (*Storage, error) {

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(endpointURL),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(strings.HasSuffix(endpointURL, "http://")),
		S3ForcePathStyle: aws.Bool(true),
	}
	awsSession := session.New(s3Config)

	return &Storage{
		awsSession: awsSession,
		bucket:     aws.String(bucket),
	}, nil
}

// NewSession creates a new session
func (s *Storage) NewSession(tenant *models.Tenant) blob.Session {
	return &Session{
		storage:  s,
		tenant:   tenant,
		s3Client: s3.New(s.awsSession),
	}
}

func (s *Session) keyFullPathURL(key string) string {
	if s.tenant != nil {
		return path.Join("tenants", strconv.Itoa(s.tenant.ID), key)
	}
	return key
}

// Get returns a blob with given key
func (s *Session) Get(key string) (*blob.Blob, error) {
	resp, err := s.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: s.storage.bucket,
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
func (s *Session) Delete(key string) error {
	_, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.storage.bucket,
		Key:    aws.String(s.keyFullPathURL(key)),
	})
	if err != nil {
		if isNotFound(err) {
			return blob.ErrNotFound
		}
		return errors.Wrap(err, "failed to delete blob '%s' from S3", key)
	}
	return nil
}

// Store a blob with given key and content. Blobs with same key are replaced.
func (s *Session) Store(b *blob.Blob) error {
	if err := blob.ValidateKey(b.Key); err != nil {
		return errors.Wrap(err, "failed to validate blob key '%s'", b.Key)
	}

	reader := bytes.NewReader(b.Object)
	_, err := s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      s.storage.bucket,
		Key:         aws.String(s.keyFullPathURL(b.Key)),
		ContentType: aws.String(b.ContentType),
		ACL:         aws.String(s3.ObjectCannedACLPrivate),
		Body:        reader,
	})
	if err != nil {
		return errors.Wrap(err, "failed to upload blob '%s' to S3", b.Key)
	}
	return nil
}
