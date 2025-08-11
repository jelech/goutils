package s3util

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseS3Path(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *S3Path
		expectError bool
	}{
		{
			name:  "valid s3 path",
			input: "s3://my-bucket/path/to/file.txt",
			expected: &S3Path{
				Bucket: "my-bucket",
				Key:    "path/to/file.txt",
			},
			expectError: false,
		},
		{
			name:  "s3 path with root object",
			input: "s3://my-bucket/file.txt",
			expected: &S3Path{
				Bucket: "my-bucket",
				Key:    "file.txt",
			},
			expectError: false,
		},
		{
			name:  "s3 path bucket only",
			input: "s3://my-bucket/",
			expected: &S3Path{
				Bucket: "my-bucket",
				Key:    "",
			},
			expectError: false,
		},
		{
			name:        "invalid scheme",
			input:       "http://my-bucket/file.txt",
			expectError: true,
		},
		{
			name:        "missing bucket",
			input:       "s3:///file.txt",
			expectError: true,
		},
		{
			name:        "empty path",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseS3Path(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Bucket, result.Bucket)
				assert.Equal(t, tt.expected.Key, result.Key)
			}
		})
	}
}

func TestS3Path_String(t *testing.T) {
	path := &S3Path{
		Bucket: "test-bucket",
		Key:    "path/to/file.txt",
	}

	expected := "s3://test-bucket/path/to/file.txt"
	assert.Equal(t, expected, path.String())
}

func TestNewClient(t *testing.T) {
	config := &Config{
		Region:           "us-east-1",
		AccessKeyID:      "test-access-key",
		SecretAccessKey:  "test-secret-key",
		Endpoint:         "http://localhost:4566", // LocalStack endpoint
		DisableSSL:       true,
		S3ForcePathStyle: true,
	}

	client, err := NewClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.s3Client)
	assert.NotNil(t, client.uploader)
	assert.NotNil(t, client.downloader)
	assert.Equal(t, "us-east-1", client.region)
}

// Mock client for testing without actual S3
type mockS3Client struct {
	objects map[string][]byte
}

func newMockS3Client() *mockS3Client {
	return &mockS3Client{
		objects: make(map[string][]byte),
	}
}

func (m *mockS3Client) putObject(bucket, key string, data []byte) {
	path := bucket + "/" + key
	m.objects[path] = data
}

func (m *mockS3Client) getObject(bucket, key string) ([]byte, bool) {
	path := bucket + "/" + key
	data, exists := m.objects[path]
	return data, exists
}

func (m *mockS3Client) deleteObject(bucket, key string) {
	path := bucket + "/" + key
	delete(m.objects, path)
}

func (m *mockS3Client) objectExists(bucket, key string) bool {
	path := bucket + "/" + key
	_, exists := m.objects[path]
	return exists
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	// Test parse and string operations
	originalPath := "s3://test-bucket/test/file.txt"
	parsed, err := ParseS3Path(originalPath)
	require.NoError(t, err)

	assert.Equal(t, "test-bucket", parsed.Bucket)
	assert.Equal(t, "test/file.txt", parsed.Key)
	assert.Equal(t, originalPath, parsed.String())
}

func TestS3PathOperations(t *testing.T) {
	testCases := []struct {
		name    string
		s3Path  string
		bucket  string
		key     string
		isValid bool
	}{
		{
			name:    "standard path",
			s3Path:  "s3://my-bucket/folder/file.txt",
			bucket:  "my-bucket",
			key:     "folder/file.txt",
			isValid: true,
		},
		{
			name:    "deep nested path",
			s3Path:  "s3://data-bucket/year=2023/month=12/day=25/data.json",
			bucket:  "data-bucket",
			key:     "year=2023/month=12/day=25/data.json",
			isValid: true,
		},
		{
			name:    "bucket root",
			s3Path:  "s3://root-bucket/",
			bucket:  "root-bucket",
			key:     "",
			isValid: true,
		},
		{
			name:    "single file",
			s3Path:  "s3://simple/file",
			bucket:  "simple",
			key:     "file",
			isValid: true,
		},
		{
			name:    "invalid protocol",
			s3Path:  "https://bucket/file",
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path, err := ParseS3Path(tc.s3Path)

			if tc.isValid {
				assert.NoError(t, err)
				assert.Equal(t, tc.bucket, path.Bucket)
				assert.Equal(t, tc.key, path.Key)

				// Test round-trip
				reconstructed := path.String()
				assert.Equal(t, tc.s3Path, reconstructed)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestFileOperations(t *testing.T) {
	// Test file operation helpers without actual S3 calls
	t.Run("test file creation and reading", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := tempDir + "/test.txt"
		testContent := []byte("hello world")

		// Write test file
		err := os.WriteFile(testFile, testContent, 0644)
		require.NoError(t, err)

		// Read test file
		content, err := os.ReadFile(testFile)
		require.NoError(t, err)
		assert.Equal(t, testContent, content)
	})

	t.Run("test bytes operations", func(t *testing.T) {
		testData := []byte("test data for s3 operations")
		reader := bytes.NewReader(testData)

		// Read all data from reader
		readData := make([]byte, len(testData))
		n, err := reader.Read(readData)
		require.NoError(t, err)
		assert.Equal(t, len(testData), n)
		assert.Equal(t, testData, readData)
	})
}

func TestS3PathValidation(t *testing.T) {
	validPaths := []string{
		"s3://bucket/file",
		"s3://my-bucket-123/path/to/file.txt",
		"s3://bucket.with.dots/file",
		"s3://bucket/",
		"s3://bucket/path/with/many/slashes/file.ext",
	}

	for _, path := range validPaths {
		t.Run("valid_"+path, func(t *testing.T) {
			parsed, err := ParseS3Path(path)
			assert.NoError(t, err)
			assert.NotNil(t, parsed)
			assert.NotEmpty(t, parsed.Bucket)

			// Test that we can reconstruct the path
			reconstructed := parsed.String()
			assert.Equal(t, path, reconstructed)
		})
	}

	invalidPaths := []string{
		"http://bucket/file",
		"s3:/file",
		"s3://",
		"bucket/file",
		"",
		"s3:///file",
	}

	for _, path := range invalidPaths {
		t.Run("invalid_"+path, func(t *testing.T) {
			parsed, err := ParseS3Path(path)
			assert.Error(t, err)
			assert.Nil(t, parsed)
		})
	}
}

func TestClientConfiguration(t *testing.T) {
	t.Run("minimal config", func(t *testing.T) {
		config := &Config{
			Region: "us-west-2",
		}

		client, err := NewClient(config)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "us-west-2", client.region)
	})

	t.Run("full config", func(t *testing.T) {
		config := &Config{
			Region:           "eu-west-1",
			AccessKeyID:      "AKIATEST",
			SecretAccessKey:  "secret",
			SessionToken:     "token",
			Endpoint:         "http://localhost:4566",
			DisableSSL:       true,
			S3ForcePathStyle: true,
		}

		client, err := NewClient(config)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "eu-west-1", client.region)
	})
}

// Benchmark tests
func BenchmarkParseS3Path(b *testing.B) {
	path := "s3://my-bucket/path/to/file.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseS3Path(path)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkS3PathString(b *testing.B) {
	path := &S3Path{
		Bucket: "my-bucket",
		Key:    "path/to/file.txt",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = path.String()
	}
}

// Integration test helpers (require actual S3 or LocalStack)
func TestIntegrationHelpers(t *testing.T) {
	// Skip if no integration test environment
	if os.Getenv("S3_INTEGRATION_TEST") == "" {
		t.Skip("Set S3_INTEGRATION_TEST=1 to run integration tests")
	}

	// Example integration test setup
	config := &Config{
		Region:           "us-east-1",
		Endpoint:         os.Getenv("S3_ENDPOINT"), // e.g., http://localhost:4566
		DisableSSL:       true,
		S3ForcePathStyle: true,
	}

	client, err := NewClient(config)
	require.NoError(t, err)

	// Test bucket and key for integration tests
	testBucket := "test-bucket-" + time.Now().Format("20060102-150405")
	testKey := "test/integration/file.txt"
	testData := []byte("integration test data")

	t.Run("put and get object", func(t *testing.T) {
		// Put object
		err := client.PutObject(testBucket, testKey, testData, "text/plain")
		if err != nil && strings.Contains(err.Error(), "NoSuchBucket") {
			t.Skip("Test bucket does not exist, skipping integration test")
		}
		require.NoError(t, err)

		// Get object
		retrievedData, err := client.GetObject(testBucket, testKey)
		require.NoError(t, err)
		assert.Equal(t, testData, retrievedData)

		// Clean up
		err = client.DeleteObject(testBucket, testKey)
		assert.NoError(t, err)
	})
}
