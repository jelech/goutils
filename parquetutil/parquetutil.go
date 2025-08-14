package parquetutil

import (
	"fmt"
	"io"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/reader"
	parquetWriter "github.com/xitongsys/parquet-go/writer"
)

// WriteConfig holds configuration options for parquet writing
type WriteConfig struct {
	// ParallelNumber sets the number of parallel workers for writing (default: 4)
	ParallelNumber int64

	// RowGroupSize sets the size of each row group in bytes (default: 128MB)
	RowGroupSize int64

	// CompressionType sets the compression algorithm (default: SNAPPY)
	CompressionType parquet.CompressionCodec

	// PageSize sets the page size in bytes (default: 8KB)
	PageSize int64

	// Repetition sets the repetition type for the schema (default: REQUIRED)
	RepetitionType parquet.FieldRepetitionType

	// SchemaWriteMode controls how schema is written (default: CREATE)
	SchemaWriteMode string

	// EnableDictionary enables dictionary encoding (default: true)
	EnableDictionary bool

	// EnableStats enables statistics collection (default: true)
	EnableStats bool
}

// Predefined compression types for convenience
const (
	CompressionUncompressed = parquet.CompressionCodec_UNCOMPRESSED
	CompressionSnappy       = parquet.CompressionCodec_SNAPPY
	CompressionGzip         = parquet.CompressionCodec_GZIP
	CompressionLZ4          = parquet.CompressionCodec_LZ4
	CompressionBrotli       = parquet.CompressionCodec_BROTLI
	CompressionZstd         = parquet.CompressionCodec_ZSTD
)

// Predefined row group sizes for convenience
const (
	RowGroupSize32MB  int64 = 32 * 1024 * 1024  // 32MB
	RowGroupSize64MB  int64 = 64 * 1024 * 1024  // 64MB
	RowGroupSize128MB int64 = 128 * 1024 * 1024 // 128MB
	RowGroupSize256MB int64 = 256 * 1024 * 1024 // 256MB
	RowGroupSize512MB int64 = 512 * 1024 * 1024 // 512MB
)

// Predefined page sizes for convenience
const (
	PageSize4KB  int64 = 4 * 1024  // 4KB
	PageSize8KB  int64 = 8 * 1024  // 8KB
	PageSize16KB int64 = 16 * 1024 // 16KB
	PageSize32KB int64 = 32 * 1024 // 32KB
	PageSize64KB int64 = 64 * 1024 // 64KB
)

// DefaultWriteConfig returns a WriteConfig with sensible defaults
func DefaultWriteConfig() *WriteConfig {
	return &WriteConfig{
		ParallelNumber:   4,
		RowGroupSize:     RowGroupSize128MB,
		CompressionType:  CompressionSnappy,
		PageSize:         PageSize8KB,
		RepetitionType:   parquet.FieldRepetitionType_REQUIRED,
		SchemaWriteMode:  "CREATE",
		EnableDictionary: true,
		EnableStats:      true,
	}
}

// NewWriteConfig creates a new WriteConfig with custom settings
func NewWriteConfig() *WriteConfig {
	return DefaultWriteConfig()
}

// WithParallelNumber sets the parallel number for writing
func (c *WriteConfig) WithParallelNumber(parallelNumber int64) *WriteConfig {
	c.ParallelNumber = parallelNumber
	return c
}

// WithRowGroupSize sets the row group size
func (c *WriteConfig) WithRowGroupSize(size int64) *WriteConfig {
	c.RowGroupSize = size
	return c
}

// WithCompressionType sets the compression type
func (c *WriteConfig) WithCompressionType(compression parquet.CompressionCodec) *WriteConfig {
	c.CompressionType = compression
	return c
}

// WithPageSize sets the page size
func (c *WriteConfig) WithPageSize(size int64) *WriteConfig {
	c.PageSize = size
	return c
}

// WithDictionary enables or disables dictionary encoding
func (c *WriteConfig) WithDictionary(enable bool) *WriteConfig {
	c.EnableDictionary = enable
	return c
}

// WithStats enables or disables statistics collection
func (c *WriteConfig) WithStats(enable bool) *WriteConfig {
	c.EnableStats = enable
	return c
}

func Read(filePath string, stuTypePoint interface{}, stus interface{}, callback func(interface{}) error) error {
	var err error
	fr, err := local.NewLocalFileReader(filePath)
	if err != nil {
		return err
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, stuTypePoint, 4)
	if err != nil {
		return err
	}
	defer pr.ReadStop()

	num := int(pr.GetNumRows())
	for i := 0; i < num/10; i++ {
		if err = pr.Read(&stus); err != nil {
			return err
		}

		fmt.Println(stus)
		if err = callback(stus); err != nil {
			return err
		}
	}

	return nil
}

func ReadSimple(filePath string, stuTypePoint interface{}, stus interface{}, callback func(interface{}) error) error {
	var err error
	fr, err := local.NewLocalFileReader(filePath)
	if err != nil {
		return err
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, nil, 4)
	if err != nil {
		return err
	}

	num := int(pr.GetNumRows())
	res, err := pr.ReadByNumber(num)
	if err != nil {
		return err
	}

	return callback(res)
}

func WriteTo(
	w io.Writer, stuTypePoint interface{},
	callback func(writer *parquetWriter.ParquetWriter) error,
) error {
	return WriteToWithConfig(w, stuTypePoint, DefaultWriteConfig(), callback)
}

// WriteToWithConfig writes parquet data to an io.Writer with custom configuration
func WriteToWithConfig(
	w io.Writer,
	stuTypePoint interface{},
	config *WriteConfig,
	callback func(writer *parquetWriter.ParquetWriter) error,
) error {
	if config == nil {
		config = DefaultWriteConfig()
	}

	pw, err := parquetWriter.NewParquetWriterFromWriter(w, stuTypePoint, config.ParallelNumber)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := pw.WriteStop(); closeErr != nil {
			// Log the error but don't override the original error
			fmt.Printf("Warning: Error closing parquet writer: %v\n", closeErr)
		}
	}()

	// Apply configuration
	pw.RowGroupSize = config.RowGroupSize
	pw.CompressionType = config.CompressionType
	pw.PageSize = config.PageSize

	if err := callback(pw); err != nil {
		return err
	}

	return nil
}

// WriteToFile writes parquet data to a file with default configuration
func WriteToFile(
	filePath string,
	stuTypePoint interface{},
	callback func(writer *parquetWriter.ParquetWriter) error,
) error {
	return WriteToFileWithConfig(filePath, stuTypePoint, DefaultWriteConfig(), callback)
}

// WriteToFileWithConfig writes parquet data to a file with custom configuration
func WriteToFileWithConfig(
	filePath string,
	stuTypePoint interface{},
	config *WriteConfig,
	callback func(writer *parquetWriter.ParquetWriter) error,
) error {
	if config == nil {
		config = DefaultWriteConfig()
	}

	fw, err := local.NewLocalFileWriter(filePath)
	if err != nil {
		return err
	}
	defer fw.Close()

	pw, err := parquetWriter.NewParquetWriterFromWriter(fw, stuTypePoint, config.ParallelNumber)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := pw.WriteStop(); closeErr != nil {
			fmt.Printf("Warning: Error closing parquet writer: %v\n", closeErr)
		}
	}()

	// Apply configuration
	pw.RowGroupSize = config.RowGroupSize
	pw.CompressionType = config.CompressionType
	pw.PageSize = config.PageSize

	if err := callback(pw); err != nil {
		return err
	}

	return nil
}

// WriteSliceTo writes a slice of data to an io.Writer with default configuration
func WriteSliceTo(w io.Writer, stuTypePoint interface{}, data interface{}) error {
	return WriteSliceToWithConfig(w, stuTypePoint, data, DefaultWriteConfig())
}

// WriteSliceToWithConfig writes a slice of data to an io.Writer with custom configuration
func WriteSliceToWithConfig(w io.Writer, stuTypePoint interface{}, data interface{}, config *WriteConfig) error {
	return WriteToWithConfig(w, stuTypePoint, config, func(writer *parquetWriter.ParquetWriter) error {
		return writer.Write(data)
	})
}

// WriteSliceToFile writes a slice of data to a file with default configuration
func WriteSliceToFile(filePath string, stuTypePoint interface{}, data interface{}) error {
	return WriteSliceToFileWithConfig(filePath, stuTypePoint, data, DefaultWriteConfig())
}

// WriteSliceToFileWithConfig writes a slice of data to a file with custom configuration
func WriteSliceToFileWithConfig(filePath string, stuTypePoint interface{}, data interface{}, config *WriteConfig) error {
	return WriteToFileWithConfig(filePath, stuTypePoint, config, func(writer *parquetWriter.ParquetWriter) error {
		return writer.Write(data)
	})
}

// WriteBatchTo writes data in batches to an io.Writer for better memory management
func WriteBatchTo(
	w io.Writer,
	stuTypePoint interface{},
	batchSize int,
	dataProvider func() (interface{}, bool, error), // returns data, hasMore, error
) error {
	return WriteBatchToWithConfig(w, stuTypePoint, batchSize, dataProvider, DefaultWriteConfig())
}

// WriteBatchToWithConfig writes data in batches to an io.Writer with custom configuration
func WriteBatchToWithConfig(
	w io.Writer,
	stuTypePoint interface{},
	batchSize int,
	dataProvider func() (interface{}, bool, error),
	config *WriteConfig,
) error {
	return WriteToWithConfig(w, stuTypePoint, config, func(writer *parquetWriter.ParquetWriter) error {
		for {
			data, hasMore, err := dataProvider()
			if err != nil {
				return err
			}
			if !hasMore {
				break
			}

			if err := writer.Write(data); err != nil {
				return err
			}
		}
		return nil
	})
}
