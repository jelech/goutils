package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDo_Success(t *testing.T) {
	callCount := 0
	err := Do(func() error {
		callCount++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestDo_RetryAndSuccess(t *testing.T) {
	callCount := 0
	err := Do(func() error {
		callCount++
		if callCount < 3 {
			return errors.New("temporary error")
		}
		return nil
	}, WithMaxAttempts(5))

	assert.NoError(t, err)
	assert.Equal(t, 3, callCount)
}

func TestDo_MaxAttemptsReached(t *testing.T) {
	callCount := 0
	err := Do(func() error {
		callCount++
		return errors.New("persistent error")
	}, WithMaxAttempts(3))

	assert.Error(t, err)
	assert.Equal(t, 3, callCount)
	assert.Contains(t, err.Error(), "retry failed after 3 attempts")
}

func TestDo_PermanentError(t *testing.T) {
	callCount := 0
	err := Do(func() error {
		callCount++
		return Permanent(errors.New("permanent error"))
	}, WithMaxAttempts(5), WithRetryIf(func(err error) bool {
		return !IsPermanent(err)
	}))

	assert.Error(t, err)
	assert.Equal(t, 1, callCount)
	assert.True(t, IsPermanent(err))
}

func TestDo_WithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	callCount := 0
	err := Do(func() error {
		callCount++
		if callCount == 2 {
			cancel() // Cancel after second attempt
		}
		return errors.New("error")
	}, WithMaxAttempts(5), WithContext(ctx), WithDelay(time.Millisecond*50))

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Equal(t, 2, callCount)
}

func TestDo_OnRetryCallback(t *testing.T) {
	var retryAttempts []int
	var retryErrors []error

	err := Do(func() error {
		if len(retryAttempts) < 2 {
			return errors.New("temporary error")
		}
		return nil
	}, WithMaxAttempts(5), WithOnRetry(func(attempt int, err error) {
		retryAttempts = append(retryAttempts, attempt)
		retryErrors = append(retryErrors, err)
	}))

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, retryAttempts)
	assert.Len(t, retryErrors, 2)
}

func TestCalculateDelay(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		attempt  int
		expected time.Duration
	}{
		{
			name: "Fixed delay",
			config: &Config{
				BaseDelay: time.Second,
				Strategy:  FixedDelay,
				Jitter:    false,
				MaxDelay:  time.Minute,
			},
			attempt:  1,
			expected: time.Second,
		},
		{
			name: "Exponential backoff",
			config: &Config{
				BaseDelay: time.Second,
				Strategy:  ExponentialBackoff,
				Jitter:    false,
				MaxDelay:  time.Minute,
			},
			attempt:  2,
			expected: time.Second * 2,
		},
		{
			name: "Linear backoff",
			config: &Config{
				BaseDelay: time.Second,
				Strategy:  LinearBackoff,
				Jitter:    false,
				MaxDelay:  time.Minute,
			},
			attempt:  3,
			expected: time.Second * 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay := calculateDelay(tt.config, tt.attempt)
			assert.Equal(t, tt.expected, delay)
		})
	}
}

func TestCalculateDelay_WithMaxDelay(t *testing.T) {
	config := &Config{
		BaseDelay: time.Second,
		Strategy:  ExponentialBackoff,
		Jitter:    false,
		MaxDelay:  time.Second * 5,
	}

	// Attempt 4 would normally give 8 seconds, but should be capped at 5
	delay := calculateDelay(config, 4)
	assert.Equal(t, time.Second*5, delay)
}

func TestCalculateDelay_WithJitter(t *testing.T) {
	config := &Config{
		BaseDelay: time.Second,
		Strategy:  FixedDelay,
		Jitter:    true,
		MaxDelay:  time.Minute,
	}

	delay := calculateDelay(config, 1)
	// With jitter, delay should be between 0.5s and 1s
	assert.True(t, delay >= time.Millisecond*500)
	assert.True(t, delay <= time.Second)
}

func TestPermanentError(t *testing.T) {
	originalErr := errors.New("original error")
	permErr := Permanent(originalErr)

	assert.True(t, IsPermanent(permErr))
	assert.Equal(t, originalErr.Error(), permErr.Error())
	assert.Equal(t, originalErr, errors.Unwrap(permErr))
}

func TestIsRetryable(t *testing.T) {
	assert.False(t, IsRetryable(nil))
	assert.True(t, IsRetryable(errors.New("some error")))
}

func BenchmarkDo_Success(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := Do(func() error {
			return nil
		})
		require.NoError(b, err)
	}
}

func BenchmarkDo_WithRetries(b *testing.B) {
	for i := 0; i < b.N; i++ {
		callCount := 0
		err := Do(func() error {
			callCount++
			if callCount < 3 {
				return errors.New("temporary error")
			}
			return nil
		}, WithMaxAttempts(5), WithDelay(time.Microsecond))
		require.NoError(b, err)
	}
}
