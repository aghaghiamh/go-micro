package utils

import (
	"fmt"
	"log"
	"math"
	"time"
)

type BackoffStrategy int

const (
	ExponentialBackoff BackoffStrategy = iota
	LinearBackoff
	ConstantBackoff
)

type BackoffConfig struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffStrategy BackoffStrategy
}

// DefaultBackoffConfig returns sensible default configuration
func DefaultBackoffConfig() BackoffConfig {
	return BackoffConfig{
		MaxRetries:      5,
		InitialDelay:    1 * time.Second,
		MaxDelay:        30 * time.Second,
		BackoffStrategy: ExponentialBackoff,
	}
}

type RetryableOperation func() error

func Backoff(op RetryableOperation, conf BackoffConfig) func() error {
	return func() error {
		for attempt := 0; ; attempt++ {
			lastErr := op()

			if lastErr == nil {
				return nil
			} else if attempt >= conf.MaxRetries {
				return fmt.Errorf("backoff mechanism failed after %d attempts, last error: %w", attempt, lastErr)
			}

			delay := calculateDelay(attempt, conf)
			log.Printf("Attempt %d failed; Retrying in %v", attempt+1, delay)

			time.Sleep(delay)
		}
	}
}

func calculateDelay(attempt int, config BackoffConfig) time.Duration {
	var delay time.Duration

	switch config.BackoffStrategy {
	case ExponentialBackoff:
		delay = time.Duration(math.Pow(2, float64(attempt))) * config.InitialDelay
	case LinearBackoff:
		delay = time.Duration(attempt+1) * config.InitialDelay
	case ConstantBackoff:
		delay = config.InitialDelay
	default:
		delay = config.InitialDelay
	}

	// Cap the delay at maximum
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}

	return delay
}
