package humanize

import (
	"strconv"
	"strings"
	"time"
)

type durationConfig struct {
	unitFlags byte
	precision int
}

const (
	Nanoseconds  byte = 1
	Microseconds byte = Nanoseconds << 1
	Milliseconds byte = Microseconds << 1
	Seconds      byte = Milliseconds << 1
	Minutes      byte = Seconds << 1
	Hours        byte = Minutes << 1
	Days         byte = Hours << 1
)

const DefaultDurationFlags = Days | Hours | Minutes | Seconds | Milliseconds
const DefaultDurationPrecision = 0

type DurationOption = func(*durationConfig)

func DurationUse(unitFlags byte) DurationOption {
	return func(config *durationConfig) {
		config.unitFlags = unitFlags
	}
}

func DurationPrecision(precision int) DurationOption {
	return func(config *durationConfig) {
		config.precision = precision
	}
}

func Duration(duration time.Duration, opts ...DurationOption) string {
	config := durationConfig{
		unitFlags: DefaultDurationFlags,
		precision: DefaultDurationPrecision,
	}
	for _, opt := range opts {
		opt(&config)
	}

	if config.unitFlags == 0 {
		return "<err:no-unit>"
	}

	builder := strings.Builder{}

	if duration < 0 {
		builder.WriteString("-")
		duration = -duration
	}

	wroteValue := false

	process := func(flag byte, unit time.Duration, suffix string) {
		if config.unitFlags&flag == 0 {
			// this flag is not enabled in the config
			return
		}

		config.unitFlags -= flag
		if duration >= unit {
			if config.unitFlags == 0 {
				value := float64(duration) / float64(unit)
				builder.WriteString(formatFloat(value, config.precision))
				builder.WriteString(suffix)
				wroteValue = true
				duration = duration % unit
			} else {
				value := duration / unit
				builder.WriteString(strconv.FormatInt(int64(value), 10))
				builder.WriteString(suffix)
				wroteValue = true
				duration = duration % unit
			}
		} else if config.unitFlags == 0 && !wroteValue {
			value := float64(duration) / float64(unit)
			builder.WriteString(formatFloat(value, config.precision))
			builder.WriteString(suffix)
		}
	}

	process(Days, 24*time.Hour, "d")
	process(Hours, time.Hour, "h")
	process(Minutes, time.Minute, "m")
	process(Seconds, time.Second, "s")
	process(Milliseconds, time.Millisecond, "ms")
	process(Microseconds, time.Microsecond, "us")
	process(Nanoseconds, time.Nanosecond, "ns")

	return builder.String()
}
