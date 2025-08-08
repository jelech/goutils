package timing

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockLogger for testing
type mockLogger struct {
	output []string
}

func (m *mockLogger) Printf(format string, v ...interface{}) {
	m.output = append(m.output, fmt.Sprintf(format, v...))
}

func TestNew(t *testing.T) {
	timer := New("test")
	assert.Equal(t, "test", timer.name)
	assert.NotNil(t, timer.logger)
}

func TestNewWithLogger(t *testing.T) {
	logger := &mockLogger{}
	timer := NewWithLogger("test", logger)
	assert.Equal(t, "test", timer.name)
	assert.Equal(t, logger, timer.logger)
}

func TestTimerStartStop(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	timer := New("test operation")
	timer.Start()

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	elapsed := timer.Stop()

	// Check that elapsed time is reasonable
	assert.True(t, elapsed >= 10*time.Millisecond)
	assert.True(t, elapsed < 100*time.Millisecond)

	// Check log output
	output := buf.String()
	assert.Contains(t, output, "[TIMING] test operation took")
}

func TestTimerStopf(t *testing.T) {
	logger := &mockLogger{}
	timer := NewWithLogger("test", logger)
	timer.Start()

	time.Sleep(5 * time.Millisecond)

	elapsed := timer.Stopf("processing %d items", 100)

	assert.True(t, elapsed >= 5*time.Millisecond)
	require.Len(t, logger.output, 1)
	assert.Contains(t, logger.output[0], "[TIMING] test: processing 100 items (took")
}

func TestMeasure(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	elapsed := Measure("test function", func() {
		time.Sleep(10 * time.Millisecond)
	})

	assert.True(t, elapsed >= 10*time.Millisecond)

	output := buf.String()
	assert.Contains(t, output, "[TIMING] test function took")
}

func TestMeasureWithResult(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	result, elapsed := MeasureWithResult("calculation", func() interface{} {
		time.Sleep(5 * time.Millisecond)
		return 42
	})

	assert.Equal(t, 42, result)
	assert.True(t, elapsed >= 5*time.Millisecond)

	output := buf.String()
	assert.Contains(t, output, "[TIMING] calculation took")
}

func TestMeasureWithError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	testErr := errors.New("test error")
	elapsed, err := MeasureWithError("error function", func() error {
		time.Sleep(5 * time.Millisecond)
		return testErr
	})

	assert.Equal(t, testErr, err)
	assert.True(t, elapsed >= 5*time.Millisecond)

	output := buf.String()
	assert.Contains(t, output, "[TIMING] error function took")
}

func TestMeasureWithResultAndError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	// Test successful case
	result, elapsed, err := MeasureWithResultAndError("success function", func() (interface{}, error) {
		time.Sleep(5 * time.Millisecond)
		return "success", nil
	})

	assert.Equal(t, "success", result)
	assert.NoError(t, err)
	assert.True(t, elapsed >= 5*time.Millisecond)

	// Test error case
	testErr := errors.New("test error")
	result2, elapsed2, err2 := MeasureWithResultAndError("error function", func() (interface{}, error) {
		time.Sleep(5 * time.Millisecond)
		return "", testErr
	})

	assert.Equal(t, "", result2)
	assert.Equal(t, testErr, err2)
	assert.True(t, elapsed2 >= 5*time.Millisecond)
}

func TestWithTimer(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	wrappedFunc := WithTimer("wrapped function", func() {
		time.Sleep(5 * time.Millisecond)
	})

	wrappedFunc()

	output := buf.String()
	assert.Contains(t, output, "[TIMING] wrapped function took")
}

func TestTrack(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	func() {
		defer Track("tracked operation")()
		time.Sleep(5 * time.Millisecond)
	}()

	output := buf.String()
	assert.Contains(t, output, "[TIMING] tracked operation took")
}

func TestTrackf(t *testing.T) {
	logger := &mockLogger{}

	func() {
		timer := NewWithLogger("tracked", logger)
		timer.Start()
		defer func() {
			timer.Stopf("processed %d records", 50)
		}()
		time.Sleep(5 * time.Millisecond)
	}()

	require.Len(t, logger.output, 1)
	assert.Contains(t, logger.output[0], "[TIMING] tracked: processed 50 records (took")
}

func TestTimerSince(t *testing.T) {
	timer := New("test").Start()
	time.Sleep(10 * time.Millisecond)

	elapsed := timer.Since()
	assert.True(t, elapsed >= 10*time.Millisecond)
}

func TestTimerLap(t *testing.T) {
	logger := &mockLogger{}
	timer := NewWithLogger("race", logger)
	timer.Start()

	time.Sleep(5 * time.Millisecond)
	elapsed1 := timer.Lap("checkpoint 1")

	time.Sleep(5 * time.Millisecond)
	elapsed2 := timer.Lap("checkpoint 2")

	assert.True(t, elapsed1 >= 5*time.Millisecond)
	assert.True(t, elapsed2 >= 10*time.Millisecond)
	assert.True(t, elapsed2 > elapsed1)

	require.Len(t, logger.output, 2)
	assert.Contains(t, logger.output[0], "[TIMING] race - checkpoint 1:")
	assert.Contains(t, logger.output[1], "[TIMING] race - checkpoint 2:")
}

func TestTimerReset(t *testing.T) {
	timer := New("test").Start()
	time.Sleep(10 * time.Millisecond)

	timer.Reset()
	time.Sleep(5 * time.Millisecond)

	elapsed := timer.Since()
	assert.True(t, elapsed >= 5*time.Millisecond)
	assert.True(t, elapsed < 10*time.Millisecond)
}

func TestTimerIsRunning(t *testing.T) {
	timer := New("test")
	assert.False(t, timer.IsRunning())

	timer.Start()
	assert.True(t, timer.IsRunning())
}

func TestComplexTimingScenario(t *testing.T) {
	logger := &mockLogger{}

	// Simulate a complex operation with multiple phases
	_, elapsed := MeasureWithResult("database operation", func() interface{} {
		// Phase 1: Connection
		timer := NewWithLogger("db-connection", logger).Start()
		time.Sleep(2 * time.Millisecond)
		timer.Lap("connected")

		// Phase 2: Query
		time.Sleep(3 * time.Millisecond)
		timer.Lap("query executed")

		// Phase 3: Processing
		time.Sleep(2 * time.Millisecond)
		timer.Stop()

		return "result"
	})

	assert.True(t, elapsed >= 7*time.Millisecond)
	assert.True(t, len(logger.output) >= 3) // connected, query executed, stop
}

// Benchmark tests
func BenchmarkTimerOverhead(b *testing.B) {
	timer := New("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timer.Start()
		timer.Stop()
	}
}

func BenchmarkMeasureOverhead(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Measure("benchmark", func() {
			// Empty function to measure overhead
		})
	}
}

// Test new statistics functionality
func TestRecorder(t *testing.T) {
	recorder := NewRecorder()

	// Record some timings
	recorder.Record("operation1", 100*time.Millisecond)
	recorder.Record("operation1", 200*time.Millisecond)
	recorder.Record("operation1", 150*time.Millisecond)

	recorder.Record("operation2", 50*time.Millisecond)
	recorder.Record("operation2", 75*time.Millisecond)

	// Test operation1 stats
	stats1, exists := recorder.Get("operation1")
	require.True(t, exists)
	assert.Equal(t, "operation1", stats1.Name)
	assert.Equal(t, int64(3), stats1.Count)
	assert.Equal(t, 450*time.Millisecond, stats1.TotalTime)
	assert.Equal(t, 150*time.Millisecond, stats1.AvgTime)
	assert.Equal(t, 100*time.Millisecond, stats1.MinTime)
	assert.Equal(t, 200*time.Millisecond, stats1.MaxTime)

	// Test operation2 stats
	stats2, exists := recorder.Get("operation2")
	require.True(t, exists)
	assert.Equal(t, "operation2", stats2.Name)
	assert.Equal(t, int64(2), stats2.Count)
	assert.Equal(t, 125*time.Millisecond, stats2.TotalTime)

	// Test GetAll
	allStats := recorder.GetAll()
	assert.Len(t, allStats, 2)
	assert.Contains(t, allStats, "operation1")
	assert.Contains(t, allStats, "operation2")

	// Test reset operation
	recorder.ResetOperation("operation1")
	_, exists = recorder.Get("operation1")
	assert.False(t, exists)

	// Test reset all
	recorder.Reset()
	allStats = recorder.GetAll()
	assert.Len(t, allStats, 0)
}

func TestGlobalRecorder(t *testing.T) {
	// Clear any previous stats
	ResetStats()

	// Test MeasureAndRecord
	elapsed := MeasureAndRecord("global_test", func() {
		time.Sleep(10 * time.Millisecond)
	})

	assert.True(t, elapsed >= 10*time.Millisecond)

	// Check stats were recorded
	stats, exists := GetStats("global_test")
	require.True(t, exists)
	assert.Equal(t, int64(1), stats.Count)
	assert.True(t, stats.TotalTime >= 10*time.Millisecond)

	// Test MeasureAndRecordWithResult
	result, elapsed2 := MeasureAndRecordWithResult("global_test_result", func() interface{} {
		time.Sleep(5 * time.Millisecond)
		return "test_result"
	})

	assert.Equal(t, "test_result", result)
	assert.True(t, elapsed2 >= 5*time.Millisecond)

	// Check global stats
	allStats := GetAllStats()
	assert.Len(t, allStats, 2)

	// Clean up
	ResetStats()
}

func TestWithRecording(t *testing.T) {
	ResetStats()

	func() {
		defer WithRecording("deferred_operation")()
		time.Sleep(5 * time.Millisecond)
	}()

	stats, exists := GetStats("deferred_operation")
	require.True(t, exists)
	assert.Equal(t, int64(1), stats.Count)
	assert.True(t, stats.TotalTime >= 5*time.Millisecond)

	ResetStats()
}

func TestConcurrentRecording(t *testing.T) {
	ResetStats()

	// Simulate 50 concurrent operations
	const goroutines = 50
	const sleepTime = 2 * time.Millisecond

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			MeasureAndRecord("concurrent_operation", func() {
				time.Sleep(sleepTime)
			})
		}()
	}

	wg.Wait()

	// Check aggregated stats
	stats, exists := GetStats("concurrent_operation")
	require.True(t, exists)
	assert.Equal(t, int64(goroutines), stats.Count)
	assert.True(t, stats.TotalTime >= time.Duration(goroutines)*sleepTime)
	assert.True(t, stats.MinTime >= sleepTime)
	assert.True(t, stats.AvgTime >= sleepTime)

	t.Logf("Concurrent stats: %s", stats.String())

	ResetStats()
}

func TestStatsString(t *testing.T) {
	stats := &Stats{
		Name:      "test_operation",
		Count:     5,
		TotalTime: 500 * time.Millisecond,
		MinTime:   50 * time.Millisecond,
		MaxTime:   150 * time.Millisecond,
		AvgTime:   100 * time.Millisecond,
	}

	str := stats.String()
	assert.Contains(t, str, "test_operation")
	assert.Contains(t, str, "count=5")
	assert.Contains(t, str, "total=500ms")
	assert.Contains(t, str, "avg=100ms")
}

func TestPrintStats(t *testing.T) {
	ResetStats()

	// Record some operations
	Record("print_test1", 100*time.Millisecond)
	Record("print_test2", 200*time.Millisecond)

	// This should not panic and should print to log
	PrintAllStats()

	ResetStats()
}
