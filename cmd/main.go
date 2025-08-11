package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jelech/goutils"
	"github.com/jelech/goutils/cacheutil"
	"github.com/jelech/goutils/convutil"
	"github.com/jelech/goutils/httputil"
	"github.com/jelech/goutils/retryutil"
	"github.com/jelech/goutils/s3util"
	"github.com/jelech/goutils/strutil"
	"github.com/jelech/goutils/timeutil"
)

func main() {
	fmt.Println(goutils.Hello())
	fmt.Println()

	// Retry example
	fmt.Println("=== Retry Examples ===")
	retryExample()
	fmt.Println()

	// HTTP client example
	fmt.Println("=== HTTP Client Examples ===")
	httpExample()
	fmt.Println()

	// Cache example
	fmt.Println("=== Cache Examples ===")
	cacheExample()
	fmt.Println()

	// String utilities example
	fmt.Println("=== String Utilities Examples ===")
	stringExample()
	fmt.Println()

	// Convert utilities example
	fmt.Println("=== Convert Utilities Examples ===")
	convertExample()
	fmt.Println()

	// Parquet file operations example (not implemented in this version)
	fmt.Println("=== Parquet File Operations Examples ===")
	fmt.Println("Parquet functionality is not implemented in this version.")
	fmt.Println("For Parquet support, upgrade to Go 1.18+ and implement with generics.")
	fmt.Println()

	// Timing examples
	fmt.Println("=== Timing Examples ===")
	timingExample()
	fmt.Println()

	// S3 examples
	fmt.Println("=== S3 Examples ===")
	s3Example()
}

func retryExample() {
	attempt := 0
	err := retryutil.Do(func() error {
		attempt++
		fmt.Printf("Attempt %d\n", attempt)
		if attempt < 3 {
			return fmt.Errorf("temporary error")
		}
		return nil
	}, retryutil.WithMaxAttempts(5), retryutil.WithDelay(time.Millisecond*100))

	if err != nil {
		log.Printf("Retry failed: %v", err)
	} else {
		fmt.Println("Retry succeeded!")
	}

	// Example with context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	err = retryutil.Do(func() error {
		time.Sleep(time.Millisecond * 200)
		return fmt.Errorf("slow operation")
	}, retryutil.WithMaxAttempts(5), retryutil.WithContext(ctx))

	if err != nil {
		fmt.Printf("Context-aware retry: %v\n", err)
	}
}

func httpExample() {
	client := httputil.NewClient(
		httputil.WithTimeout(time.Second*30),
		httputil.WithHeaders(map[string]string{
			"User-Agent": "GoUtils-Example/1.0",
		}),
	)

	// Note: This example uses httpbin.org for demonstration
	// In a real application, you would use your actual API endpoints
	fmt.Println("Creating HTTP client with custom headers and timeout")
	fmt.Printf("Client configured with 30s timeout\n")

	// Simulate a request structure
	type ExampleResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	// Example of how you might use the client
	fmt.Println("HTTP client ready for API calls")
	_ = client // Prevent unused variable warning
}

func cacheExample() {
	// Memory cache example
	memCache := cacheutil.NewMemoryCache()

	// Set some values
	memCache.Set("user:1", "John Doe", time.Minute*5)
	memCache.Set("user:2", "Jane Smith", time.Minute*5)

	// Get values
	if value, exists := memCache.Get("user:1"); exists {
		fmt.Printf("Memory Cache - user:1 = %s\n", value)
	}

	fmt.Printf("Memory Cache size: %d\n", memCache.Size())

	// LRU cache example
	lruCache := cacheutil.NewLRUCache(3)

	lruCache.Set("key1", "value1", 0)
	lruCache.Set("key2", "value2", 0)
	lruCache.Set("key3", "value3", 0)

	fmt.Printf("LRU Cache size: %d\n", lruCache.Size())

	// This will evict key1 (least recently used)
	lruCache.Set("key4", "value4", 0)

	if _, exists := lruCache.Get("key1"); !exists {
		fmt.Println("key1 was evicted from LRU cache")
	}
}

func stringExample() {
	text := "hello_world_example"

	fmt.Printf("Original: %s\n", text)
	fmt.Printf("CamelCase: %s\n", strutil.CamelCase(text))
	fmt.Printf("PascalCase: %s\n", strutil.PascalCase(text))
	fmt.Printf("KebabCase: %s\n", strutil.KebabCase(text))
	fmt.Printf("Reversed: %s\n", strutil.Reverse(text))

	// Random string generation
	if randomStr, err := strutil.RandomString(10); err == nil {
		fmt.Printf("Random string: %s\n", randomStr)
	}

	// Email validation
	emails := []string{"test@example.com", "invalid.email", "user@domain.co.uk"}
	for _, email := range emails {
		valid := strutil.IsValidEmail(email)
		fmt.Printf("%s is valid email: %t\n", email, valid)
	}

	// String padding
	fmt.Printf("Padded left: '%s'\n", strutil.PadLeft("123", 8, '0'))
	fmt.Printf("Padded right: '%s'\n", strutil.PadRight("123", 8, ' '))
}

func convertExample() {
	// Type conversions
	fmt.Printf("String to int: %s -> %d\n", "42", convutil.MustToInt("42"))
	fmt.Printf("Int to string: %d -> %s\n", 42, convutil.ToString(42))
	fmt.Printf("String to bool: %s -> %t\n", "true", convutil.MustToBool("true"))

	// Slice conversions
	intSlice := []int{1, 2, 3, 4, 5}
	if strSlice, err := convutil.ToStringSlice(intSlice); err == nil {
		fmt.Printf("Int slice to string slice: %v -> %v\n", intSlice, strSlice)
	}

	// Struct to map
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		City string `json:"city"`
	}

	person := Person{Name: "John", Age: 30, City: "New York"}
	if personMap, err := convutil.ToMap(person); err == nil {
		fmt.Printf("Struct to map: %+v\n", personMap)
	}

	// JSON conversion
	if jsonStr, err := convutil.ToJSON(person); err == nil {
		fmt.Printf("Struct to JSON: %s\n", jsonStr)
	}

	// Time conversion
	timeStr := "2023-12-25T15:30:00Z"
	if t, err := convutil.ToTime(timeStr); err == nil {
		fmt.Printf("String to time: %s -> %s\n", timeStr, t.Format("2006-01-02 15:04:05"))
	}
}

// timingExample demonstrates timing package features
func timingExample() {
	// Example 1: Simple timing
	timer1 := timeutil.New("operation-1")
	timer1.Start()
	time.Sleep(100 * time.Millisecond)
	duration1 := timer1.Stop()
	fmt.Printf("Operation 1 took: %v\n", duration1)

	// Example 2: Measure function execution
	duration := timeutil.Measure("calculation", func() {
		sum := 0
		for i := 0; i < 1000000; i++ {
			sum += i
		}
	})
	fmt.Printf("Calculation took: %v\n", duration)

	// Example 3: Measure with result
	result, duration2 := timeutil.MeasureWithResult("calculation-with-result", func() interface{} {
		sum := 0
		for i := 0; i < 1000000; i++ {
			sum += i
		}
		return sum
	})
	fmt.Printf("Calculation result: %v, took: %v\n", result, duration2)

	// Example 4: Measure and record statistics
	for i := 0; i < 5; i++ {
		timeutil.MeasureAndRecord("repeated-operation", func() {
			time.Sleep(time.Duration(10+i*5) * time.Millisecond)
		})
	}

	// Example 5: Get statistics
	stats, exists := timeutil.GetStats("repeated-operation")
	if exists {
		fmt.Printf("Repeated operation stats: %s\n", stats.String())
	}

	// Example 6: Get all statistics
	fmt.Println("\nAll timing statistics:")
	allStats := timeutil.GetAllStats()
	for name, stat := range allStats {
		fmt.Printf("  %s: %s\n", name, stat.String())
	}

	// Example 7: Concurrent operations with global recorder
	fmt.Println("\nRunning concurrent operations...")

	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer func() { done <- true }()

			timeutil.MeasureAndRecord(fmt.Sprintf("concurrent-op-%d", id), func() {
				time.Sleep(time.Duration(50+id*10) * time.Millisecond)
			})
		}(i)
	}

	// Wait for all operations to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// Show updated statistics
	fmt.Println("\nConcurrent operation statistics:")
	timeutil.PrintAllStats()
} // s3Example demonstrates S3 package features
func s3Example() {
	// Note: This example assumes AWS credentials are configured
	// or uses LocalStack for local development

	config := &s3util.Config{
		Region: "us-east-1",
		// Uncomment below for LocalStack
		// Endpoint:         "http://localhost:4566",
		// DisableSSL:       true,
		// S3ForcePathStyle: true,
	}

	client, err := s3util.NewClient(config)
	if err != nil {
		fmt.Printf("Failed to create S3 client (this is expected if AWS is not configured): %v\n", err)
		fmt.Println("To use S3 examples, configure AWS credentials or use LocalStack")
		return
	}

	bucket := "goutils-example-bucket"
	key := "examples/test.txt"
	content := "Hello from GoUtils S3 package!"

	fmt.Println("Demonstrating S3 operations...")

	// Example 1: Upload string
	fmt.Println("1. Uploading string data...")
	_, err = client.UploadString(bucket, key, content, &s3util.UploadOptions{
		ContentType: "text/plain",
		Metadata: map[string]*string{
			"Source":  &[]string{"GoUtils"}[0],
			"Example": &[]string{"true"}[0],
		},
	})
	if err != nil {
		fmt.Printf("Upload failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Uploaded to s3://%s/%s\n", bucket, key)

	// Example 2: Download string
	fmt.Println("2. Downloading string data...")
	downloaded, err := client.DownloadString(bucket, key, nil)
	if err != nil {
		fmt.Printf("Download failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Downloaded: %s\n", downloaded)

	// Example 3: Check existence
	fmt.Println("3. Checking object existence...")
	exists, err := client.ObjectExists(bucket, key)
	if err != nil {
		fmt.Printf("Existence check failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Object exists: %v\n", exists)

	// Example 4: List objects
	fmt.Println("4. Listing objects...")
	objects, err := client.ListObjects(bucket, "examples/", 10)
	if err != nil {
		fmt.Printf("List failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Found %d objects\n", len(objects))

	// Example 5: Path operations
	fmt.Println("5. S3 path operations...")
	s3Path := "s3://example-bucket/folder/file.txt"

	parsed, err := s3util.ParseS3Path(s3Path)
	if err != nil {
		fmt.Printf("Path parsing failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Parsed path - Bucket: %s, Key: %s\n", parsed.Bucket, parsed.Key)

	dir, _ := s3util.GetS3PathDir(s3Path)
	base, _ := s3util.GetS3PathBase(s3Path)
	fmt.Printf("✓ Directory: %s, Base: %s\n", dir, base)

	fmt.Println("S3 operations completed successfully!")
}
