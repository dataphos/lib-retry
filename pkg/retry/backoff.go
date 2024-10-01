// Copyright 2024 Syntio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import (
	"time"
)

// Backoff models a backoff strategy.
type Backoff interface {
	// Next returns the time duration to wait and whether to stop.
	Next() (time.Duration, bool)
}

// BackoffFunc is a backoff expressed as a function.
type BackoffFunc func() (time.Duration, bool)

// Next implements Backoff.
func (b BackoffFunc) Next() (time.Duration, bool) {
	return b()
}

// WithMaxRetries decorates the given Backoff, setting a maximum amount of retries.
func WithMaxRetries(max int, next Backoff) Backoff {
	if max <= 0 {
		return noRetry
	}

	var attempt int

	return BackoffFunc(func() (time.Duration, bool) {
		if attempt >= max {
			return 0, true
		}
		attempt++

		val, stop := next.Next()
		if stop {
			return 0, true
		}

		return val, false
	})
}

var noRetry = BackoffFunc(func() (time.Duration, bool) {
	return 0, true
})

// Exponential returns an exponential backoff implementation of Backoff.
func Exponential(base time.Duration) Backoff {
	var attempt int

	return BackoffFunc(func() (time.Duration, bool) {
		attempt++

		return base << (attempt - 1), false
	})
}

// Constant returns a constant backoff implementation of Backoff.
func Constant(interval time.Duration) Backoff {
	return BackoffFunc(func() (time.Duration, bool) {
		return interval, false
	})
}
