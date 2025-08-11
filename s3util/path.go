package s3util

import (
	"fmt"
	"net/url"
	"strings"
)

// S3Path represents a parsed S3 path
type S3Path struct {
	Bucket string
	Key    string
	Region string
}

// String returns the S3 path as a URL string
func (p *S3Path) String() string {
	return fmt.Sprintf("s3://%s/%s", p.Bucket, p.Key)
}

// ParseS3Path parses an S3 URL into bucket and key components
func ParseS3Path(s3Path string) (*S3Path, error) {
	if !strings.HasPrefix(s3Path, "s3://") {
		return nil, fmt.Errorf("invalid S3 path: %s (must start with s3://)", s3Path)
	}

	// Parse URL
	u, err := url.Parse(s3Path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse S3 path: %w", err)
	}

	bucket := u.Host
	key := strings.TrimPrefix(u.Path, "/")

	if bucket == "" {
		return nil, fmt.Errorf("bucket name is required in S3 path: %s", s3Path)
	}

	return &S3Path{
		Bucket: bucket,
		Key:    key,
	}, nil
}

// JoinS3Path creates an S3 path from bucket and key
func JoinS3Path(bucket, key string) string {
	return fmt.Sprintf("s3://%s/%s", bucket, key)
}

// SplitS3Path splits an S3 path into bucket and key
func SplitS3Path(s3Path string) (bucket, key string, err error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return "", "", err
	}
	return path.Bucket, path.Key, nil
}

// IsValidS3Path checks if a string is a valid S3 path
func IsValidS3Path(s3Path string) bool {
	_, err := ParseS3Path(s3Path)
	return err == nil
}

// GetS3PathDir returns the directory part of an S3 path
func GetS3PathDir(s3Path string) (string, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return "", err
	}

	if path.Key == "" {
		return s3Path, nil
	}

	lastSlash := strings.LastIndex(path.Key, "/")
	if lastSlash == -1 {
		return fmt.Sprintf("s3://%s/", path.Bucket), nil
	}

	return fmt.Sprintf("s3://%s/%s", path.Bucket, path.Key[:lastSlash+1]), nil
}

// GetS3PathBase returns the base name of an S3 path
func GetS3PathBase(s3Path string) (string, error) {
	path, err := ParseS3Path(s3Path)
	if err != nil {
		return "", err
	}

	if path.Key == "" {
		return "", nil
	}

	lastSlash := strings.LastIndex(path.Key, "/")
	if lastSlash == -1 {
		return path.Key, nil
	}

	return path.Key[lastSlash+1:], nil
}
