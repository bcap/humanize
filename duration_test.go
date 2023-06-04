package humanize

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	var useFlags byte
	var precision int

	test := func(value time.Duration, expected string, opts ...DurationOption) {
		opts = append([]DurationOption{DurationUse(useFlags), DurationPrecision(precision)}, opts...)
		assert.Equal(t, expected, Duration(value, opts...))
		isZero, err := regexp.MatchString(`^0\w{1,2}$`, expected)
		if err != nil {
			panic(err)
		}
		if !isZero {
			assert.Equal(t, "-"+expected, Duration(-value, opts...))
		}
	}

	useFlags = Days | Hours | Minutes | Seconds | Milliseconds | Microseconds | Nanoseconds
	precision = 0

	test(time.Nanosecond, "1ns")
	test(time.Microsecond, "1us")
	test(time.Millisecond, "1ms")
	test(time.Second, "1s")
	test(time.Minute, "1m")
	test(time.Hour, "1h")
	test(24*time.Hour, "1d")
	test(10*24*time.Hour, "10d")

	test(15*time.Second, "15s")
	test(48*time.Minute+15*time.Second, "48m15s")
	test(3*time.Hour+48*time.Minute+15*time.Second, "3h48m15s")
	test(8*24*time.Hour+3*time.Hour+48*time.Minute+15*time.Second, "8d3h48m15s")
	test(8*24*time.Hour+3*time.Hour+48*time.Minute+15*time.Second+380*time.Millisecond, "8d3h48m15s380ms")
	test(8*24*time.Hour+3*time.Hour+48*time.Minute+15*time.Second+380*time.Millisecond+639*time.Microsecond, "8d3h48m15s380ms639us")
	test(8*24*time.Hour+3*time.Hour+48*time.Minute+15*time.Second+380*time.Millisecond+639*time.Microsecond+714*time.Nanosecond, "8d3h48m15s380ms639us714ns")

	dur := 8*24*time.Hour + 3*time.Hour + 48*time.Minute + 15*time.Second + 380*time.Millisecond + 639*time.Microsecond + 714*time.Nanosecond
	test(dur, "8d", DurationUse(Days))
	test(dur, "8d4h", DurationUse(Days|Hours))
	test(dur, "8d3h48m", DurationUse(Days|Hours|Minutes))
	test(dur, "8d3h48m15s", DurationUse(Days|Hours|Minutes|Seconds))
	test(dur, "8d3h48m15s381ms", DurationUse(Days|Hours|Minutes|Seconds|Milliseconds))
	test(dur, "8d3h48m15s380ms640us", DurationUse(Days|Hours|Minutes|Seconds|Milliseconds|Microseconds))
	test(dur, "8d3h48m15s380ms639us714ns", DurationUse(Days|Hours|Minutes|Seconds|Milliseconds|Microseconds|Nanoseconds))

	useFlags = Milliseconds
	precision = 3

	test(time.Nanosecond, "0ms")
	test(time.Microsecond, "0.001ms")
	test(time.Millisecond, "1ms")
	test(time.Second, "1000ms")
	test(time.Minute, "60000ms")
	test(time.Hour, "3600000ms")

	useFlags = Seconds | Milliseconds
	precision = 3

	test(1500*time.Millisecond, "1s500ms")
	test(1500*time.Millisecond+100*time.Microsecond, "1s500.1ms")
	test(1500*time.Millisecond+10*time.Microsecond, "1s500.01ms")
	test(1500*time.Millisecond+time.Microsecond, "1s500.001ms")
	test(1500*time.Millisecond+100*time.Nanosecond, "1s500ms")
}
