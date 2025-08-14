package parquetutil

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xitongsys/parquet-go/parquet"
	parquetWriter "github.com/xitongsys/parquet-go/writer"
)

// TestData represents a simple test structure
type TestData struct {
	ID   int64  `parquet:"name=id, type=INT64"`
	Name string `parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8"`
	Age  int32  `parquet:"name=age, type=INT32"`
}

func TestDefaultWriteConfig(t *testing.T) {
	config := DefaultWriteConfig()

	assert.Equal(t, int64(4), config.ParallelNumber)
	assert.Equal(t, RowGroupSize128MB, config.RowGroupSize)
	assert.Equal(t, CompressionSnappy, config.CompressionType)
	assert.Equal(t, PageSize8KB, config.PageSize)
	assert.True(t, config.EnableDictionary)
	assert.True(t, config.EnableStats)
}

func TestWriteConfigBuilder(t *testing.T) {
	config := NewWriteConfig().
		WithParallelNumber(8).
		WithRowGroupSize(RowGroupSize64MB).
		WithCompressionType(CompressionGzip).
		WithPageSize(PageSize16KB).
		WithDictionary(false).
		WithStats(false)

	assert.Equal(t, int64(8), config.ParallelNumber)
	assert.Equal(t, RowGroupSize64MB, config.RowGroupSize)
	assert.Equal(t, CompressionGzip, config.CompressionType)
	assert.Equal(t, PageSize16KB, config.PageSize)
	assert.False(t, config.EnableDictionary)
	assert.False(t, config.EnableStats)
}

func TestWriteSliceTo(t *testing.T) {
	data := []TestData{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	var buf bytes.Buffer
	err := WriteSliceTo(&buf, new(TestData), data)
	assert.NoError(t, err)
	assert.Greater(t, buf.Len(), 0)
}

func TestWriteSliceToWithConfig(t *testing.T) {
	data := []TestData{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
	}

	config := NewWriteConfig().
		WithCompressionType(CompressionGzip).
		WithRowGroupSize(RowGroupSize32MB)

	var buf bytes.Buffer
	err := WriteSliceToWithConfig(&buf, new(TestData), data, config)
	assert.NoError(t, err)
	assert.Greater(t, buf.Len(), 0)
}

func TestWriteTo(t *testing.T) {
	var buf bytes.Buffer

	err := WriteTo(&buf, new(TestData), func(writer *parquetWriter.ParquetWriter) error {
		data := []TestData{
			{ID: 1, Name: "Test", Age: 20},
		}
		return writer.Write(data)
	})

	assert.NoError(t, err)
	assert.Greater(t, buf.Len(), 0)
}

func TestWriteToWithConfig(t *testing.T) {
	var buf bytes.Buffer

	config := NewWriteConfig().
		WithCompressionType(CompressionZstd).
		WithParallelNumber(2)

	err := WriteToWithConfig(&buf, new(TestData), config, func(writer *parquetWriter.ParquetWriter) error {
		data := []TestData{
			{ID: 1, Name: "Test", Age: 20},
		}
		return writer.Write(data)
	})

	assert.NoError(t, err)
	assert.Greater(t, buf.Len(), 0)
}

func TestCompressionConstants(t *testing.T) {
	assert.Equal(t, parquet.CompressionCodec_UNCOMPRESSED, CompressionUncompressed)
	assert.Equal(t, parquet.CompressionCodec_SNAPPY, CompressionSnappy)
	assert.Equal(t, parquet.CompressionCodec_GZIP, CompressionGzip)
	assert.Equal(t, parquet.CompressionCodec_LZ4, CompressionLZ4)
	assert.Equal(t, parquet.CompressionCodec_BROTLI, CompressionBrotli)
	assert.Equal(t, parquet.CompressionCodec_ZSTD, CompressionZstd)
}

func TestSizeConstants(t *testing.T) {
	// Test row group sizes
	assert.Equal(t, int64(32*1024*1024), RowGroupSize32MB)
	assert.Equal(t, int64(64*1024*1024), RowGroupSize64MB)
	assert.Equal(t, int64(128*1024*1024), RowGroupSize128MB)
	assert.Equal(t, int64(256*1024*1024), RowGroupSize256MB)
	assert.Equal(t, int64(512*1024*1024), RowGroupSize512MB)

	// Test page sizes
	assert.Equal(t, int64(4*1024), PageSize4KB)
	assert.Equal(t, int64(8*1024), PageSize8KB)
	assert.Equal(t, int64(16*1024), PageSize16KB)
	assert.Equal(t, int64(32*1024), PageSize32KB)
	assert.Equal(t, int64(64*1024), PageSize64KB)
}

type Data struct {
	Idx          string    `parquet:"name=idx, type=BYTE_ARRAY, convertedtype=UTF8"`
	Date         int64     `parquet:"name=date, type=INT64"`
	Actual       int64     `parquet:"name=actual, type=INT64"`
	PredY        float64   `parquet:"name=pred_y, type=DOUBLE"`
	Demand       []float64 `parquet:"name=demand, type=DOUBLE, repetitiontype=REPEATED"`
	InitialStock float64   `parquet:"name=initial_stock, type=DOUBLE"`
}

func TestRead(t *testing.T) {
	// Create a temporary parquet file for testing
	tempFile := "test_data.parquet"

	// First write some test data
	err := WriteToFile(tempFile, &Data{}, func(writer *parquetWriter.ParquetWriter) error {
		for i := 0; i < 10; i++ {
			testData := &Data{
				Idx:          fmt.Sprintf("test_%d", i),
				Date:         int64(i),
				Actual:       int64(i * 3),
				PredY:        float64(i) * 1.5,
				Demand:       []float64{float64(i), float64(i + 1)},
				InitialStock: float64(i) * 2.5,
			}
			if err := writer.Write(testData); err != nil {
				return fmt.Errorf("failed to write test data: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
		return
	}

	// Clean up temp file after test
	defer func() {
		if err := os.Remove(tempFile); err != nil {
			t.Logf("Warning: failed to remove temp file %s: %v", tempFile, err)
		}
	}()

	// Now test reading
	typeData := Data{}
	arrData := make([]Data, 10)
	err = ReadSimple(tempFile, &typeData, arrData, func(i interface{}) error {
		// ReadSimple passes the raw data from parquet reader
		// The type is not necessarily []Data but the raw parquet data
		if i == nil {
			return fmt.Errorf("received nil data")
		}
		// Just verify we got some data
		return nil
	})
	if err != nil {
		t.Errorf("Read error: %v", err)
		return
	}
}

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	err := WriteTo(&buf, &Data{}, func(writer *parquetWriter.ParquetWriter) error {
		for i := 0; i < 10; i++ {
			if err := writer.Write(&Data{}); err != nil {
				return fmt.Errorf("failed to write schema: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		t.Errorf("Write error: %v", err)
	}
}
