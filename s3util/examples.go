package s3util

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"
)

// Examples demonstrates various S3 operations
func Examples() {
	// Configure S3 client
	config := &Config{
		Region:           "us-east-1",
		Endpoint:         "http://localhost:4566", // LocalStack endpoint for testing
		DisableSSL:       true,
		S3ForcePathStyle: true,
	}

	client, err := NewClient(config)
	if err != nil {
		log.Printf("Failed to create S3 client: %v", err)
		return
	}

	bucket := "example-bucket"
	key := "test/example.txt"
	content := "Hello, S3 World!"

	fmt.Println("=== S3 Operations Examples ===")

	// Example 1: Upload string data
	fmt.Println("\n1. Upload string data:")
	uploadOptions := &UploadOptions{
		ContentType: "text/plain",
		Metadata: map[string]*string{
			"Author":  &[]string{"GoUtils"}[0],
			"Version": &[]string{"1.0.0"}[0],
		},
	}

	_, err = client.UploadString(bucket, key, content, uploadOptions)
	if err != nil {
		log.Printf("Upload failed: %v", err)
	} else {
		fmt.Printf("✓ Uploaded string to s3://%s/%s\n", bucket, key)
	}

	// Example 2: Download as string
	fmt.Println("\n2. Download as string:")
	downloadedContent, err := client.DownloadString(bucket, key, nil)
	if err != nil {
		log.Printf("Download failed: %v", err)
	} else {
		fmt.Printf("✓ Downloaded content: %s\n", downloadedContent)
	}

	// Example 3: Upload bytes
	fmt.Println("\n3. Upload bytes:")
	data := []byte("Binary data example")
	binaryKey := "test/binary.dat"

	_, err = client.UploadBytes(bucket, binaryKey, data, &UploadOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		log.Printf("Binary upload failed: %v", err)
	} else {
		fmt.Printf("✓ Uploaded binary data to s3://%s/%s\n", bucket, binaryKey)
	}

	// Example 4: Download bytes
	fmt.Println("\n4. Download bytes:")
	downloadedData, err := client.DownloadBytes(bucket, binaryKey, nil)
	if err != nil {
		log.Printf("Binary download failed: %v", err)
	} else {
		fmt.Printf("✓ Downloaded %d bytes\n", len(downloadedData))
	}

	// Example 5: Upload from reader
	fmt.Println("\n5. Upload from reader:")
	reader := strings.NewReader("Data from reader")
	readerKey := "test/from-reader.txt"

	_, err = client.UploadFromReader(bucket, readerKey, reader, &UploadOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		log.Printf("Reader upload failed: %v", err)
	} else {
		fmt.Printf("✓ Uploaded from reader to s3://%s/%s\n", bucket, readerKey)
	}

	// Example 6: Use S3 path strings
	fmt.Println("\n6. Using S3 path strings:")
	s3Path := "s3://example-bucket/test/path-example.txt"

	_, err = client.UploadStringToPath(s3Path, "Content via S3 path", &UploadOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		log.Printf("Path upload failed: %v", err)
	} else {
		fmt.Printf("✓ Uploaded to %s\n", s3Path)
	}

	// Example 7: Check if object exists
	fmt.Println("\n7. Check object existence:")
	exists, err := client.ObjectExists(bucket, key)
	if err != nil {
		log.Printf("Existence check failed: %v", err)
	} else {
		fmt.Printf("✓ Object s3://%s/%s exists: %v\n", bucket, key, exists)
	}

	// Example 8: Generate presigned URL
	fmt.Println("\n8. Generate presigned URL:")
	presignedURL, err := client.GetPresignedURL(bucket, key, 1*time.Hour)
	if err != nil {
		log.Printf("Presigned URL generation failed: %v", err)
	} else {
		fmt.Printf("✓ Presigned URL: %s\n", presignedURL[:50]+"...")
	}

	// Example 9: List objects
	fmt.Println("\n9. List objects:")
	objects, err := client.ListObjects(bucket, "test/", 10)
	if err != nil {
		log.Printf("List objects failed: %v", err)
	} else {
		fmt.Printf("✓ Found %d objects with prefix 'test/'\n", len(objects))
		for _, obj := range objects {
			fmt.Printf("  - %s (%d bytes)\n", *obj.Key, *obj.Size)
		}
	}

	// Example 10: Stream upload with progress
	fmt.Println("\n10. Stream upload with progress:")
	largeContent := strings.Repeat("Large content data ", 1000)
	progressKey := "test/progress-example.txt"

	progressFn := func(written, total int64) {
		if total > 0 {
			percent := float64(written) / float64(total) * 100
			fmt.Printf("\rProgress: %.1f%% (%d/%d bytes)", percent, written, total)
		}
	}

	_, err = client.StreamUpload(bucket, progressKey, strings.NewReader(largeContent),
		int64(len(largeContent)), &UploadOptions{ContentType: "text/plain"}, progressFn)
	if err != nil {
		log.Printf("Stream upload failed: %v", err)
	} else {
		fmt.Printf("\n✓ Stream upload completed\n")
	}

	// Example 11: Download to buffer
	fmt.Println("\n11. Download to buffer:")
	buffer, err := client.DownloadToBuffer(bucket, key, nil)
	if err != nil {
		log.Printf("Download to buffer failed: %v", err)
	} else {
		fmt.Printf("✓ Downloaded to buffer: %d bytes\n", buffer.Len())
	}

	// Example 12: Path operations
	fmt.Println("\n12. S3 Path operations:")
	testPath := "s3://my-bucket/folder/subfolder/file.txt"

	parsed, err := ParseS3Path(testPath)
	if err != nil {
		log.Printf("Path parsing failed: %v", err)
	} else {
		fmt.Printf("✓ Parsed path: bucket=%s, key=%s\n", parsed.Bucket, parsed.Key)
	}

	dir, err := GetS3PathDir(testPath)
	if err != nil {
		log.Printf("Get dir failed: %v", err)
	} else {
		fmt.Printf("✓ Directory: %s\n", dir)
	}

	base, err := GetS3PathBase(testPath)
	if err != nil {
		log.Printf("Get base failed: %v", err)
	} else {
		fmt.Printf("✓ Base name: %s\n", base)
	}

	fmt.Println("\n=== S3 Examples completed ===")
}

// ExampleMultipartUpload demonstrates multipart upload for large files
func ExampleMultipartUpload() {
	config := &Config{
		Region: "us-east-1",
	}

	client, err := NewClient(config)
	if err != nil {
		log.Printf("Failed to create S3 client: %v", err)
		return
	}

	// Create a large buffer to simulate a big file
	largeData := bytes.Repeat([]byte("This is test data for multipart upload. "), 100000)
	reader := bytes.NewReader(largeData)

	bucket := "large-files-bucket"
	key := "large-file.dat"

	fmt.Println("=== Multipart Upload Example ===")

	// Upload with 5MB parts
	partSize := int64(5 * 1024 * 1024) // 5MB
	concurrency := 3

	_, err = client.ConcurrentUpload(bucket, key, reader, partSize, concurrency, &UploadOptions{
		ContentType: "application/octet-stream",
		Metadata: map[string]*string{
			"UploadMethod": &[]string{"multipart"}[0],
		},
	})

	if err != nil {
		log.Printf("Multipart upload failed: %v", err)
	} else {
		fmt.Printf("✓ Multipart upload completed for %d bytes\n", len(largeData))
	}
}

// ExampleBatchOperations demonstrates batch operations
func ExampleBatchOperations() {
	config := &Config{
		Region: "us-east-1",
	}

	client, err := NewClient(config)
	if err != nil {
		log.Printf("Failed to create S3 client: %v", err)
		return
	}

	bucket := "batch-operations-bucket"
	fmt.Println("=== Batch Operations Example ===")

	// Upload multiple files
	files := map[string]string{
		"file1.txt": "Content of file 1",
		"file2.txt": "Content of file 2",
		"file3.txt": "Content of file 3",
	}

	fmt.Println("Uploading multiple files:")
	for filename, content := range files {
		key := "batch/" + filename
		_, err := client.UploadString(bucket, key, content, &UploadOptions{
			ContentType: "text/plain",
		})
		if err != nil {
			log.Printf("Failed to upload %s: %v", filename, err)
		} else {
			fmt.Printf("✓ Uploaded %s\n", filename)
		}
	}

	// List and download all files
	fmt.Println("\nListing and downloading files:")
	objects, err := client.ListObjects(bucket, "batch/", 100)
	if err != nil {
		log.Printf("Failed to list objects: %v", err)
		return
	}

	for _, obj := range objects {
		content, err := client.DownloadString(bucket, *obj.Key, nil)
		if err != nil {
			log.Printf("Failed to download %s: %v", *obj.Key, err)
		} else {
			fmt.Printf("✓ Downloaded %s: %s\n", *obj.Key, content)
		}
	}
}
