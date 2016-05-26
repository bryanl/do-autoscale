package backoff

import (
	"fmt"
	"math/rand"
	"time"
)

// Policy implements a backoff policy.
type Policy struct {
	Millis []int
}

// DefaultPolicy is a backoff policy ranging up to 10 seconds.
var DefaultPolicy = Policy{
	[]int{0, 10, 10, 100, 100, 500, 500, 3000, 300, 5000, 10000},
}

// Duration returns the time duration of the n'th wait cycle in a
// backoff policy. This is b.Millis[n] randomized to avoid
// thundering herds.
func (p Policy) Duration(n int) (time.Duration, error) {
	if n >= len(p.Millis) {
		return 0, fmt.Errorf("backoff policy exausted")
	}

	return time.Duration(jitter(p.Millis[n])) * time.Millisecond, nil
}

func jitter(millis int) int {
	if millis == 0 {
		return 0
	}

	return millis/2 + rand.Intn(millis)
}
