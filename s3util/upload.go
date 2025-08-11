package s3util

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadOptions contains options for upload operations
type UploadOptions struct {
	ContentType          string
	ContentEncoding      string
	Metadata             map[string]*string
	ACL                  string
	StorageClass         string
	ServerSideEncryption string
	KMSKeyID             string
}

// UploadFromReader uploads data from an io.Reader to S3
func (c *Client) UploadFromReader(bucket, key string, reader io.Reader, options *UploadOptions) (*s3manager.UploadOutput, error) {
	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	if options != nil {
		if options.ContentType != "" {
			input.ContentType = aws.String(options.ContentType)
		}
		if options.ContentEncoding != "" {
			input.ContentEncoding = aws.String(options.ContentEncoding)
		}
		if options.Metadata != nil {
			input.Metadata = options.Metadata
		}
		if options.ACL != "" {
			input.ACL = aws.String(options.ACL)
		}
		if options.StorageClass != "" {
			input.StorageClass = aws.String(options.StorageClass)
		}
		if options.ServerSideEncryption != "" {
			input.ServerSideEncryption = aws.String(options.ServerSideEncryption)
		}
		if options.KMSKeyID != "" {
			input.SSEKMSKeyId = aws.String(options.KMSKeyID)
		}
	}

	result, err := c.uploader.Upload(input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to s3://%s/%s: %w", bucket, key, err)
	}

	return result, nil
}

// UploadFromReaderToPath uploads data from an io.Reader using S3 path string
func (c *Client) UploadFromReaderToPath(s3Path string, reader io.Reader, options *UploadOptions) (*s3manager.UploadOutput, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return nil, err
	}

	return c.UploadFromReader(path.Bucket, path.Key, reader, options)
}

// UploadBytes uploads byte data to S3
func (c *Client) UploadBytes(bucket, key string, data []byte, options *UploadOptions) (*s3manager.UploadOutput, error) {
	reader := bytes.NewReader(data)
	return c.UploadFromReader(bucket, key, reader, options)
}

// UploadBytesToPath uploads byte data using S3 path string
func (c *Client) UploadBytesToPath(s3Path string, data []byte, options *UploadOptions) (*s3manager.UploadOutput, error) {
	reader := bytes.NewReader(data)
	return c.UploadFromReaderToPath(s3Path, reader, options)
}

// UploadString uploads string data to S3
func (c *Client) UploadString(bucket, key string, data string, options *UploadOptions) (*s3manager.UploadOutput, error) {
	reader := bytes.NewReader([]byte(data))
	return c.UploadFromReader(bucket, key, reader, options)
}

// UploadStringToPath uploads string data using S3 path string
func (c *Client) UploadStringToPath(s3Path string, data string, options *UploadOptions) (*s3manager.UploadOutput, error) {
	reader := bytes.NewReader([]byte(data))
	return c.UploadFromReaderToPath(s3Path, reader, options)
}

// UploadFile uploads a file to S3
func (c *Client) UploadFile(bucket, key, filename string, options *UploadOptions) (*s3manager.UploadOutput, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	return c.UploadFromReader(bucket, key, file, options)
}

// UploadFileToPath uploads a file using S3 path string
func (c *Client) UploadFileToPath(filename, s3Path string, options *UploadOptions) (*s3manager.UploadOutput, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return nil, err
	}

	return c.UploadFile(path.Bucket, path.Key, filename, options)
}

// StreamUpload uploads data from an io.Reader with progress callback
func (c *Client) StreamUpload(bucket, key string, reader io.Reader, size int64, options *UploadOptions, progressFn func(written, total int64)) (*s3manager.UploadOutput, error) {
	var progressReader io.Reader = reader

	if progressFn != nil && size > 0 {
		progressReader = &progressReaderWrapper{
			Reader:   reader,
			total:    size,
			callback: progressFn,
		}
	}

	return c.UploadFromReader(bucket, key, progressReader, options)
}

// progressReaderWrapper wraps an io.Reader to track progress
type progressReaderWrapper struct {
	io.Reader
	written  int64
	total    int64
	callback func(written, total int64)
}

func (pr *progressReaderWrapper) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.written += int64(n)

	if pr.callback != nil {
		pr.callback(pr.written, pr.total)
	}

	return n, err
}

// MultipartUpload performs a multipart upload for large files
func (c *Client) MultipartUpload(bucket, key string, reader io.Reader, partSize int64, options *UploadOptions) (*s3manager.UploadOutput, error) {
	uploader := s3manager.NewUploaderWithClient(c.s3Client, func(u *s3manager.Uploader) {
		if partSize > 0 {
			u.PartSize = partSize
		}
	})

	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	if options != nil {
		if options.ContentType != "" {
			input.ContentType = aws.String(options.ContentType)
		}
		if options.ContentEncoding != "" {
			input.ContentEncoding = aws.String(options.ContentEncoding)
		}
		if options.Metadata != nil {
			input.Metadata = options.Metadata
		}
		if options.ACL != "" {
			input.ACL = aws.String(options.ACL)
		}
		if options.StorageClass != "" {
			input.StorageClass = aws.String(options.StorageClass)
		}
		if options.ServerSideEncryption != "" {
			input.ServerSideEncryption = aws.String(options.ServerSideEncryption)
		}
		if options.KMSKeyID != "" {
			input.SSEKMSKeyId = aws.String(options.KMSKeyID)
		}
	}

	result, err := uploader.Upload(input)
	if err != nil {
		return nil, fmt.Errorf("failed to multipart upload to s3://%s/%s: %w", bucket, key, err)
	}

	return result, nil
}

// ConcurrentUpload uploads multiple parts concurrently
func (c *Client) ConcurrentUpload(bucket, key string, reader io.Reader, partSize int64, concurrency int, options *UploadOptions) (*s3manager.UploadOutput, error) {
	uploader := s3manager.NewUploaderWithClient(c.s3Client, func(u *s3manager.Uploader) {
		if partSize > 0 {
			u.PartSize = partSize
		}
		if concurrency > 0 {
			u.Concurrency = concurrency
		}
	})

	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	if options != nil {
		if options.ContentType != "" {
			input.ContentType = aws.String(options.ContentType)
		}
		if options.ContentEncoding != "" {
			input.ContentEncoding = aws.String(options.ContentEncoding)
		}
		if options.Metadata != nil {
			input.Metadata = options.Metadata
		}
		if options.ACL != "" {
			input.ACL = aws.String(options.ACL)
		}
		if options.StorageClass != "" {
			input.StorageClass = aws.String(options.StorageClass)
		}
		if options.ServerSideEncryption != "" {
			input.ServerSideEncryption = aws.String(options.ServerSideEncryption)
		}
		if options.KMSKeyID != "" {
			input.SSEKMSKeyId = aws.String(options.KMSKeyID)
		}
	}

	result, err := uploader.Upload(input)
	if err != nil {
		return nil, fmt.Errorf("failed to concurrent upload to s3://%s/%s: %w", bucket, key, err)
	}

	return result, nil
}
