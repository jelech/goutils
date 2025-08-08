package timer

import (
	"bytes"
	"errors"
	"fmt"
	"log"
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
	assert.Contains(t, output, "[TIMER] test operation took")
}

func TestTimerStopf(t *testing.T) {
	logger := &mockLogger{}
	timer := NewWithLogger("test", logger)
	timer.Start()

	time.Sleep(5 * time.Millisecond)

	elapsed := timer.Stopf("processing %d items", 100)

	assert.True(t, elapsed >= 5*time.Millisecond)
	require.Len(t, logger.output, 1)
	assert.Contains(t, logger.output[0], "[TIMER] test: processing 100 items (took")
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
	assert.Contains(t, output, "[TIMER] test function took")
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
	assert.Contains(t, output, "[TIMER] calculation took")
}

func TestMeasureWithError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	testErr := errors.New("test error")
	err, elapsed := MeasureWithError("error function", func() error {
		time.Sleep(5 * time.Millisecond)
		return testErr
	})

	assert.Equal(t, testErr, err)
	assert.True(t, elapsed >= 5*time.Millisecond)

	output := buf.String()
	assert.Contains(t, output, "[TIMER] error function took")
}

func TestMeasureWithResultAndError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	// Test successful case
	result, err, elapsed := MeasureWithResultAndError("success function", func() (interface{}, error) {
		time.Sleep(5 * time.Millisecond)
		return "success", nil
	})

	assert.Equal(t, "success", result)
	assert.NoError(t, err)
	assert.True(t, elapsed >= 5*time.Millisecond)

	// Test error case
	testErr := errors.New("test error")
	result2, err2, elapsed2 := MeasureWithResultAndError("error function", func() (interface{}, error) {
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
	assert.Contains(t, output, "[TIMER] wrapped function took")
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
	assert.Contains(t, output, "[TIMER] tracked operation took")
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
	assert.Contains(t, logger.output[0], "[TIMER] tracked: processed 50 records (took")
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
	assert.Contains(t, logger.output[0], "[TIMER] race - checkpoint 1:")
	assert.Contains(t, logger.output[1], "[TIMER] race - checkpoint 2:")
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
