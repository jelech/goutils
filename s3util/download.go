package s3util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// DownloadOptions contains options for download operations
type DownloadOptions struct {
	Range             string // e.g., "bytes=0-1023"
	VersionId         string
	IfMatch           string
	IfNoneMatch       string
	IfModifiedSince   *time.Time
	IfUnmodifiedSince *time.Time
}

// DownloadToWriter downloads an S3 object to an io.Writer
func (c *Client) DownloadToWriter(bucket, key string, writer io.WriterAt, options *DownloadOptions) (int64, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if options != nil {
		if options.Range != "" {
			input.Range = aws.String(options.Range)
		}
		if options.VersionId != "" {
			input.VersionId = aws.String(options.VersionId)
		}
		if options.IfMatch != "" {
			input.IfMatch = aws.String(options.IfMatch)
		}
		if options.IfNoneMatch != "" {
			input.IfNoneMatch = aws.String(options.IfNoneMatch)
		}
		if options.IfModifiedSince != nil {
			input.IfModifiedSince = options.IfModifiedSince
		}
		if options.IfUnmodifiedSince != nil {
			input.IfUnmodifiedSince = options.IfUnmodifiedSince
		}
	}

	numBytes, err := c.downloader.Download(writer, input)
	if err != nil {
		return 0, fmt.Errorf("failed to download s3://%s/%s: %w", bucket, key, err)
	}

	return numBytes, nil
}

// DownloadToWriterFromPath downloads using S3 path string to a writer
func (c *Client) DownloadToWriterFromPath(s3Path string, writer io.WriterAt, options *DownloadOptions) (int64, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return 0, err
	}

	return c.DownloadToWriter(path.Bucket, path.Key, writer, options)
}

// DownloadToBuffer downloads an S3 object to a bytes.Buffer
func (c *Client) DownloadToBuffer(bucket, key string, options *DownloadOptions) (*bytes.Buffer, error) {
	buf := &aws.WriteAtBuffer{}

	_, err := c.DownloadToWriter(bucket, key, buf, options)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(buf.Bytes()), nil
}

// DownloadToBufferFromPath downloads using S3 path string to a buffer
func (c *Client) DownloadToBufferFromPath(s3Path string, options *DownloadOptions) (*bytes.Buffer, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return nil, err
	}

	return c.DownloadToBuffer(path.Bucket, path.Key, options)
}

// DownloadBytes downloads an S3 object and returns its content as bytes
func (c *Client) DownloadBytes(bucket, key string, options *DownloadOptions) ([]byte, error) {
	buf, err := c.DownloadToBuffer(bucket, key, options)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DownloadBytesFromPath downloads using S3 path string and returns bytes
func (c *Client) DownloadBytesFromPath(s3Path string, options *DownloadOptions) ([]byte, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return nil, err
	}

	return c.DownloadBytes(path.Bucket, path.Key, options)
}

// DownloadString downloads an S3 object and returns its content as string
func (c *Client) DownloadString(bucket, key string, options *DownloadOptions) (string, error) {
	data, err := c.DownloadBytes(bucket, key, options)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// DownloadStringFromPath downloads using S3 path string and returns string
func (c *Client) DownloadStringFromPath(s3Path string, options *DownloadOptions) (string, error) {
	data, err := c.DownloadBytesFromPath(s3Path, options)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// DownloadFileAdvanced downloads an S3 object directly to a file with options
func (c *Client) DownloadFileAdvanced(bucket, key, filename string, options *DownloadOptions) error {
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
	_, err = c.DownloadToWriter(bucket, key, file, options)
	if err != nil {
		return fmt.Errorf("failed to download s3://%s/%s to %s: %w", bucket, key, filename, err)
	}

	return nil
}

// DownloadFileAdvancedFromPath downloads using S3 path string with options
func (c *Client) DownloadFileAdvancedFromPath(s3Path, filename string, options *DownloadOptions) error {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return err
	}

	return c.DownloadFileAdvanced(path.Bucket, path.Key, filename, options)
}

// StreamDownload downloads with progress callback
func (c *Client) StreamDownload(bucket, key string, writer io.WriterAt, progressFn func(written, total int64)) (int64, error) {
	var progressWriter io.WriterAt = writer

	if progressFn != nil {
		// Get object size first
		headOutput, err := c.s3Client.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return 0, fmt.Errorf("failed to get object size: %w", err)
		}

		size := *headOutput.ContentLength
		progressWriter = &progressWriterWrapper{
			WriterAt: writer,
			total:    size,
			callback: progressFn,
		}
	}

	return c.DownloadToWriter(bucket, key, progressWriter, nil)
}

// progressWriterWrapper wraps an io.WriterAt to track progress
type progressWriterWrapper struct {
	io.WriterAt
	written  int64
	total    int64
	callback func(written, total int64)
}

func (pw *progressWriterWrapper) WriteAt(p []byte, off int64) (int, error) {
	n, err := pw.WriterAt.WriteAt(p, off)
	pw.written += int64(n)

	if pw.callback != nil {
		pw.callback(pw.written, pw.total)
	}

	return n, err
}

// ConcurrentDownload downloads using multiple concurrent parts
func (c *Client) ConcurrentDownload(bucket, key string, writer io.WriterAt, partSize int64, concurrency int) (int64, error) {
	downloader := s3manager.NewDownloaderWithClient(c.s3Client, func(d *s3manager.Downloader) {
		if partSize > 0 {
			d.PartSize = partSize
		}
		if concurrency > 0 {
			d.Concurrency = concurrency
		}
	})

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	numBytes, err := downloader.Download(writer, input)
	if err != nil {
		return 0, fmt.Errorf("failed to concurrent download s3://%s/%s: %w", bucket, key, err)
	}

	return numBytes, nil
}

// PartialDownload downloads a specific range of bytes
func (c *Client) PartialDownload(bucket, key string, start, end int64, writer io.WriterAt) (int64, error) {
	options := &DownloadOptions{
		Range: fmt.Sprintf("bytes=%d-%d", start, end),
	}

	return c.DownloadToWriter(bucket, key, writer, options)
}
