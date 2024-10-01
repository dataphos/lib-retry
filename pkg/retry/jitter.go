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
	"math/rand"
	"time"
)

// WithJitter wraps Backoff by adding random jitter between 0 and 1000 milliseconds.
//
// This prevents a large amounts of clients from retrying at the same time.
func WithJitter(next Backoff) Backoff {
	return BackoffFunc(func() (time.Duration, bool) {
		v, stop := next.Next()
		if stop {
			return 0, true
		}
		// #nosec G404; math/rand is fine for this use case.
		return v + time.Duration(rand.Intn(1000))*time.Millisecond, false
	})
}
