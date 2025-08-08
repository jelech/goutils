package timer

import (
	"fmt"
	"log"
	"time"
)

// Timer represents a timer for measuring execution time
type Timer struct {
	name      string
	startTime time.Time
	logger    Logger
}

// Logger interface for custom logging implementations
type Logger interface {
	Printf(format string, v ...interface{})
}

// defaultLogger uses the standard log package
type defaultLogger struct{}

func (d defaultLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// New creates a new timer with the given name
func New(name string) *Timer {
	return &Timer{
		name:   name,
		logger: defaultLogger{},
	}
}

// NewWithLogger creates a new timer with a custom logger
func NewWithLogger(name string, logger Logger) *Timer {
	return &Timer{
		name:   name,
		logger: logger,
	}
}

// Start starts the timer
func (t *Timer) Start() *Timer {
	t.startTime = time.Now()
	return t
}

// Stop stops the timer and prints the elapsed time
func (t *Timer) Stop() time.Duration {
	elapsed := time.Since(t.startTime)
	t.logger.Printf("[TIMER] %s took %v", t.name, elapsed)
	return elapsed
}

// Stopf stops the timer and prints the elapsed time with custom format
func (t *Timer) Stopf(format string, args ...interface{}) time.Duration {
	elapsed := time.Since(t.startTime)
	message := fmt.Sprintf(format, args...)
	t.logger.Printf("[TIMER] %s: %s (took %v)", t.name, message, elapsed)
	return elapsed
}

// Measure is a convenience function to measure a function's execution time
func Measure(name string, fn func()) time.Duration {
	timer := New(name).Start()
	fn()
	return timer.Stop()
}

// MeasureWithResult measures a function that returns a value
func MeasureWithResult(name string, fn func() interface{}) (interface{}, time.Duration) {
	timer := New(name).Start()
	result := fn()
	elapsed := timer.Stop()
	return result, elapsed
}

// MeasureWithError measures a function that returns an error
func MeasureWithError(name string, fn func() error) (error, time.Duration) {
	timer := New(name).Start()
	err := fn()
	elapsed := timer.Stop()
	return err, elapsed
}

// MeasureWithResultAndError measures a function that returns both value and error
func MeasureWithResultAndError(name string, fn func() (interface{}, error)) (interface{}, error, time.Duration) {
	timer := New(name).Start()
	result, err := fn()
	elapsed := timer.Stop()
	return result, err, elapsed
}

// WithTimer is a decorator function that wraps any function with timing
func WithTimer(name string, fn func()) func() {
	return func() {
		Measure(name, fn)
	}
}

// Track creates a timer that automatically stops when the returned function is called
// Useful with defer statement
func Track(name string) func() {
	timer := New(name).Start()
	return func() {
		timer.Stop()
	}
}

// Trackf creates a timer with custom message formatting
func Trackf(name string, format string, args ...interface{}) func() {
	timer := New(name).Start()
	return func() {
		timer.Stopf(format, args...)
	}
}

// Since returns the elapsed time since the timer started
func (t *Timer) Since() time.Duration {
	return time.Since(t.startTime)
}

// Lap records a lap time without stopping the timer
func (t *Timer) Lap(lapName string) time.Duration {
	elapsed := time.Since(t.startTime)
	t.logger.Printf("[TIMER] %s - %s: %v", t.name, lapName, elapsed)
	return elapsed
}

// Reset resets the timer to current time
func (t *Timer) Reset() *Timer {
	t.startTime = time.Now()
	return t
}

// IsRunning returns true if the timer has been started
func (t *Timer) IsRunning() bool {
	return !t.startTime.IsZero()
}
