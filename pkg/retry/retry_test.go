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

package retry_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/dataphos/lib-retry/pkg/retry"
)

func TestDo(t *testing.T) {
	t.Run("exit_on_max_attempt", func(t *testing.T) {
		ctx := context.Background()
		backoff := retry.WithMaxRetries(3, retry.BackoffFunc(func() (time.Duration, bool) {
			return 1 * time.Nanosecond, false
		}))

		var counter int
		if err := retry.Do(ctx, backoff, func(_ context.Context) error {
			counter++

			return errors.New("something's wrong")
		}); err == nil {
			t.Error("expected err")
		}

		if counter != 4 {
			t.Errorf("expected %v to be %v", counter, 4)
		}
	})

	t.Run("exit_no_error", func(t *testing.T) {
		ctx := context.Background()
		backoff := retry.WithMaxRetries(3, retry.BackoffFunc(func() (time.Duration, bool) {
			return 1 * time.Millisecond, false
		}))

		var counter int
		if err := retry.Do(ctx, backoff, func(_ context.Context) error {
			counter++

			return nil
		}); err != nil {
			t.Fatal("expected no err")
		}

		if got, want := counter, 1; got != want {
			t.Errorf("expected %v to be %v", got, want)
		}
	})

	t.Run("context_canceled", func(t *testing.T) {
		t.Parallel()

		backoffFunc := retry.BackoffFunc(func() (time.Duration, bool) {
			return 5 * time.Second, false
		})

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		if err := retry.Do(ctx, backoffFunc, func(_ context.Context) error {
			return errors.New("something's wrong")
		}); !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected error to be %v, got %v", context.DeadlineExceeded, err)
		}
	})
}
