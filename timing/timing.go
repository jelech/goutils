package timing

import (
	"fmt"
	"log"
	"sync"
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
	if log.Writer() != nil {
		log.Printf(format, v...)
	}
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
	t.logger.Printf("[TIMING] %s took %v", t.name, elapsed)
	return elapsed
}

// Stopf stops the timer and prints the elapsed time with custom format
func (t *Timer) Stopf(format string, args ...interface{}) time.Duration {
	elapsed := time.Since(t.startTime)
	message := fmt.Sprintf(format, args...)
	t.logger.Printf("[TIMING] %s: %s (took %v)", t.name, message, elapsed)
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
func MeasureWithError(name string, fn func() error) (time.Duration, error) {
	timer := New(name).Start()
	err := fn()
	elapsed := timer.Stop()
	return elapsed, err
}

// MeasureWithResultAndError measures a function that returns both value and error
func MeasureWithResultAndError(name string, fn func() (interface{}, error)) (interface{}, time.Duration, error) {
	timer := New(name).Start()
	result, err := fn()
	elapsed := timer.Stop()
	return result, elapsed, err
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
	t.logger.Printf("[TIMING] %s - %s: %v", t.name, lapName, elapsed)
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

// Stats represents timing statistics for a named operation
type Stats struct {
	Name        string        `json:"name"`
	Count       int64         `json:"count"`
	TotalTime   time.Duration `json:"total_time"`
	MinTime     time.Duration `json:"min_time"`
	MaxTime     time.Duration `json:"max_time"`
	AvgTime     time.Duration `json:"avg_time"`
	LastUpdated time.Time     `json:"last_updated"`
}

// String returns a formatted string representation of the stats
func (s *Stats) String() string {
	return fmt.Sprintf("[STATS] %s: count=%d, total=%v, avg=%v, min=%v, max=%v",
		s.Name, s.Count, s.TotalTime, s.AvgTime, s.MinTime, s.MaxTime)
}

// Recorder manages timing statistics for multiple operations
type Recorder struct {
	mu    sync.RWMutex
	stats map[string]*Stats
}

// NewRecorder creates a new timing recorder
func NewRecorder() *Recorder {
	return &Recorder{
		stats: make(map[string]*Stats),
	}
}

// Record records a timing measurement
func (r *Recorder) Record(name string, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stat, exists := r.stats[name]
	if !exists {
		stat = &Stats{
			Name:    name,
			MinTime: duration,
			MaxTime: duration,
		}
		r.stats[name] = stat
	}

	stat.Count++
	stat.TotalTime += duration
	stat.AvgTime = time.Duration(int64(stat.TotalTime) / stat.Count)
	stat.LastUpdated = time.Now()

	if duration < stat.MinTime {
		stat.MinTime = duration
	}
	if duration > stat.MaxTime {
		stat.MaxTime = duration
	}
}

// Get returns the statistics for a named operation
func (r *Recorder) Get(name string) (*Stats, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stat, exists := r.stats[name]
	if !exists {
		return nil, false
	}

	// Return a copy to avoid concurrent access issues
	statsCopy := *stat
	return &statsCopy, true
}

// GetAll returns all recorded statistics
func (r *Recorder) GetAll() map[string]*Stats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]*Stats)
	for name, stat := range r.stats {
		statsCopy := *stat
		result[name] = &statsCopy
	}
	return result
}

// Reset clears all statistics
func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stats = make(map[string]*Stats)
}

// ResetOperation clears statistics for a specific operation
func (r *Recorder) ResetOperation(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.stats, name)
}

// PrintStats prints all statistics
func (r *Recorder) PrintStats() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.stats) == 0 {
		if log.Writer() != nil {
			log.Println("[TIMING] No statistics recorded")
		}
		return
	}

	if log.Writer() != nil {
		log.Println("[TIMING] Performance Statistics:")
		for _, stat := range r.stats {
			log.Println(stat.String())
		}
	}
}

// Global recorder instance
var globalRecorder = NewRecorder()

// Record records timing to the global recorder
func Record(name string, duration time.Duration) {
	globalRecorder.Record(name, duration)
}

// GetStats returns statistics from the global recorder
func GetStats(name string) (*Stats, bool) {
	return globalRecorder.Get(name)
}

// GetAllStats returns all statistics from the global recorder
func GetAllStats() map[string]*Stats {
	return globalRecorder.GetAll()
}

// ResetStats clears all global statistics
func ResetStats() {
	globalRecorder.Reset()
}

// PrintAllStats prints all global statistics
func PrintAllStats() {
	globalRecorder.PrintStats()
}

// MeasureAndRecord measures execution time and records it to global stats
func MeasureAndRecord(name string, fn func()) time.Duration {
	timer := New(name).Start()
	fn()
	elapsed := timer.Stop()
	Record(name, elapsed)
	return elapsed
}

// MeasureAndRecordWithResult measures and records with return value
func MeasureAndRecordWithResult(name string, fn func() interface{}) (interface{}, time.Duration) {
	timer := New(name).Start()
	result := fn()
	elapsed := timer.Stop()
	Record(name, elapsed)
	return result, elapsed
}

// MeasureAndRecordWithError measures and records with error return
func MeasureAndRecordWithError(name string, fn func() error) (time.Duration, error) {
	timer := New(name).Start()
	err := fn()
	elapsed := timer.Stop()
	Record(name, elapsed)
	return elapsed, err
}

// WithRecording creates a timer that automatically records to global stats
func WithRecording(name string) func() {
	timer := New(name).Start()
	return func() {
		elapsed := timer.Stop()
		Record(name, elapsed)
	}
}
