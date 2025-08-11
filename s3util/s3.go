package s3util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Client wraps AWS S3 client with convenient methods
type Client struct {
	s3Client   *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	session    *session.Session
	region     string
}

// Config holds S3 client configuration
type Config struct {
	Region           string
	AccessKeyID      string
	SecretAccessKey  string
	SessionToken     string
	Endpoint         string
	DisableSSL       bool
	S3ForcePathStyle bool
}

// NewClient creates a new S3 client with the given configuration
func NewClient(config *Config) (*Client, error) {
	awsConfig := &aws.Config{
		Region: aws.String(config.Region),
	}

	if config.Endpoint != "" {
		awsConfig.Endpoint = aws.String(config.Endpoint)
	}

	if config.DisableSSL {
		awsConfig.DisableSSL = aws.Bool(true)
	}

	if config.S3ForcePathStyle {
		awsConfig.S3ForcePathStyle = aws.Bool(true)
	}

	// Create session
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create S3 client and upload/download managers
	s3Client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)

	return &Client{
		s3Client:   s3Client,
		uploader:   uploader,
		downloader: downloader,
		session:    sess,
		region:     config.Region,
	}, nil
}

// NewClientFromSession creates a new S3 client from an existing AWS session
func NewClientFromSession(sess *session.Session) *Client {
	return &Client{
		s3Client:   s3.New(sess),
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
		session:    sess,
		region:     *sess.Config.Region,
	}
}

// GetObject downloads an object from S3 and returns its content as bytes
func (c *Client) GetObject(bucket, key string) ([]byte, error) {
	result, err := c.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object s3://%s/%s: %w", bucket, key, err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}

	return data, nil
}

// GetObjectFromPath downloads an object using S3 path string
func (c *Client) GetObjectFromPath(s3Path string) ([]byte, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return nil, err
	}

	return c.GetObject(path.Bucket, path.Key)
}

// PutObject uploads data to S3
func (c *Client) PutObject(bucket, key string, data []byte, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	_, err := c.s3Client.PutObject(input)
	if err != nil {
		return fmt.Errorf("failed to put object s3://%s/%s: %w", bucket, key, err)
	}

	return nil
}

// PutObjectFromPath uploads data using S3 path string
func (c *Client) PutObjectFromPath(s3Path string, data []byte, contentType string) error {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return err
	}

	return c.PutObject(path.Bucket, path.Key, data, contentType)
}

// DownloadFile downloads an S3 object directly to a file
func (c *Client) DownloadFile(bucket, key, filename string) error {
	// Create the directories in the path if they don't exist
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	// Download the file
	_, err = c.downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to download s3://%s/%s to %s: %w", bucket, key, filename, err)
	}

	return nil
}

// DownloadFileFromPath downloads using S3 path string
func (c *Client) DownloadFileFromPath(s3Path, filename string) error {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return err
	}

	return c.DownloadFile(path.Bucket, path.Key, filename)
}

// ObjectExists checks if an object exists in S3
func (c *Client) ObjectExists(bucket, key string) (bool, error) {
	_, err := c.s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}

// ObjectExistsFromPath checks if an object exists using S3 path string
func (c *Client) ObjectExistsFromPath(s3Path string) (bool, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return false, err
	}

	return c.ObjectExists(path.Bucket, path.Key)
}

// DeleteObject deletes an object from S3
func (c *Client) DeleteObject(bucket, key string) error {
	_, err := c.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete s3://%s/%s: %w", bucket, key, err)
	}

	return nil
}

// DeleteObjectFromPath deletes an object using S3 path string
func (c *Client) DeleteObjectFromPath(s3Path string) error {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return err
	}

	return c.DeleteObject(path.Bucket, path.Key)
}

// ListObjects lists objects in a bucket with optional prefix
func (c *Client) ListObjects(bucket, prefix string, maxKeys int64) ([]*s3.Object, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	if maxKeys > 0 {
		input.MaxKeys = aws.Int64(maxKeys)
	}

	result, err := c.s3Client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in s3://%s: %w", bucket, err)
	}

	return result.Contents, nil
}

// GetPresignedURL generates a presigned URL for an S3 object
func (c *Client) GetPresignedURL(bucket, key string, expiration time.Duration) (string, error) {
	req, _ := c.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}

// GetPresignedURLFromPath generates a presigned URL using S3 path string
func (c *Client) GetPresignedURLFromPath(s3Path string, expiration time.Duration) (string, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return "", err
	}

	return c.GetPresignedURL(path.Bucket, path.Key, expiration)
}

// CopyObject copies an object within S3
func (c *Client) CopyObject(sourceBucket, sourceKey, destBucket, destKey string) error {
	copySource := fmt.Sprintf("%s/%s", sourceBucket, sourceKey)

	_, err := c.s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(destBucket),
		Key:        aws.String(destKey),
		CopySource: aws.String(copySource),
	})
	if err != nil {
		return fmt.Errorf("failed to copy s3://%s/%s to s3://%s/%s: %w",
			sourceBucket, sourceKey, destBucket, destKey, err)
	}

	return nil
}

// GetS3Client returns the underlying S3 client for advanced operations
func (c *Client) GetS3Client() *s3.S3 {
	return c.s3Client
}

// GetUploader returns the S3 uploader for advanced upload operations
func (c *Client) GetUploader() *s3manager.Uploader {
	return c.uploader
}

// GetDownloader returns the S3 downloader for advanced download operations
func (c *Client) GetDownloader() *s3manager.Downloader {
	return c.downloader
}

// GetSession returns the underlying AWS session
func (c *Client) GetSession() *session.Session {
	return c.session
}
