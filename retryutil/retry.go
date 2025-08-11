// Package retry provides a flexible retry mechanism with configurable strategies.
package retryutil

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryableFunc represents a function that can be retried
type RetryableFunc func() error

// Strategy defines the retry strategy
type Strategy int

const (
	// FixedDelay uses a fixed delay between retries
	FixedDelay Strategy = iota
	// ExponentialBackoff uses exponential backoff with optional jitter
	ExponentialBackoff
	// LinearBackoff uses linear backoff
	LinearBackoff
)

// Config holds the configuration for retry operations
type Config struct {
	MaxAttempts int                          // Maximum number of attempts (including the first one)
	BaseDelay   time.Duration                // Base delay between retries
	MaxDelay    time.Duration                // Maximum delay between retries
	Strategy    Strategy                     // Retry strategy
	Jitter      bool                         // Whether to add jitter to delays
	RetryIf     func(error) bool             // Function to determine if an error should trigger a retry
	OnRetry     func(attempt int, err error) // Callback function called on each retry
	Context     context.Context              // Context for cancellation
}

// Option represents a configuration option for retry
type Option func(*Config)

// WithMaxAttempts sets the maximum number of attempts
func WithMaxAttempts(attempts int) Option {
	return func(c *Config) {
		c.MaxAttempts = attempts
	}
}

// WithDelay sets the base delay between retries
func WithDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.BaseDelay = delay
	}
}

// WithMaxDelay sets the maximum delay between retries
func WithMaxDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.MaxDelay = delay
	}
}

// WithBackoff sets the retry strategy
func WithBackoff(strategy Strategy) Option {
	return func(c *Config) {
		c.Strategy = strategy
	}
}

// WithJitter enables or disables jitter
func WithJitter(enabled bool) Option {
	return func(c *Config) {
		c.Jitter = enabled
	}
}

// WithRetryIf sets a custom function to determine if an error should trigger a retry
func WithRetryIf(fn func(error) bool) Option {
	return func(c *Config) {
		c.RetryIf = fn
	}
}

// WithOnRetry sets a callback function that is called on each retry
func WithOnRetry(fn func(attempt int, err error)) Option {
	return func(c *Config) {
		c.OnRetry = fn
	}
}

// WithContext sets the context for cancellation
func WithContext(ctx context.Context) Option {
	return func(c *Config) {
		c.Context = ctx
	}
}

// defaultConfig returns the default retry configuration
func defaultConfig() *Config {
	return &Config{
		MaxAttempts: 3,
		BaseDelay:   time.Millisecond * 100,
		MaxDelay:    time.Second * 30,
		Strategy:    ExponentialBackoff,
		Jitter:      true,
		RetryIf:     func(error) bool { return true },
		OnRetry:     func(int, error) {},
		Context:     context.Background(),
	}
}

// Do executes the given function with retry logic
func Do(fn RetryableFunc, options ...Option) error {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}

	var lastErr error
	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-config.Context.Done():
			return config.Context.Err()
		default:
		}

		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should retry this error
		if !config.RetryIf(err) {
			return err
		}

		// Don't wait after the last attempt
		if attempt == config.MaxAttempts {
			break
		}

		// Call the retry callback
		config.OnRetry(attempt, err)

		// Calculate delay
		delay := calculateDelay(config, attempt)

		// Wait for the delay or context cancellation
		timer := time.NewTimer(delay)
		select {
		case <-config.Context.Done():
			timer.Stop()
			return config.Context.Err()
		case <-timer.C:
		}
	}

	return fmt.Errorf("retry failed after %d attempts, last error: %w", config.MaxAttempts, lastErr)
}

// calculateDelay calculates the delay for the next retry based on the strategy
func calculateDelay(config *Config, attempt int) time.Duration {
	var delay time.Duration

	switch config.Strategy {
	case FixedDelay:
		delay = config.BaseDelay
	case ExponentialBackoff:
		delay = time.Duration(float64(config.BaseDelay) * math.Pow(2, float64(attempt-1)))
	case LinearBackoff:
		delay = time.Duration(int64(config.BaseDelay) * int64(attempt))
	default:
		delay = config.BaseDelay
	}

	// Apply jitter if enabled
	if config.Jitter {
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay = delay/2 + jitter
	}

	// Ensure delay doesn't exceed max delay
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}

	return delay
}

// IsRetryable checks if an error is retryable based on common patterns
func IsRetryable(err error) bool {
	// Add common retryable error patterns here
	// For example, network timeouts, temporary failures, etc.
	return err != nil
}

// IsTemporary checks if an error implements the Temporary interface
func IsTemporary(err error) bool {
	type temporary interface {
		Temporary() bool
	}

	if t, ok := err.(temporary); ok {
		return t.Temporary()
	}

	return false
}

// PermanentError wraps an error to indicate it should not be retried
type PermanentError struct {
	Err error
}

func (e PermanentError) Error() string {
	return e.Err.Error()
}

func (e PermanentError) Unwrap() error {
	return e.Err
}

// Permanent wraps an error to indicate it should not be retried
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return PermanentError{Err: err}
}

// IsPermanent checks if an error is permanent (should not be retried)
func IsPermanent(err error) bool {
	var permErr PermanentError
	return errors.As(err, &permErr)
}
