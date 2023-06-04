package humanize

import (
	"math"
	"strconv"
)

const (
	kib float64 = 1024
	mib float64 = kib * 1024
	gib float64 = mib * 1024
	tib float64 = gib * 1024
	pib float64 = tib * 1024
)

type bytesConfig struct {
	precision          int
	roundingMode       roundingMode
	traillingZero      bool
	valueUnitSeparator string
}

type BytesOption = func(*bytesConfig)

func TruncateBytes(config *bytesConfig) {
	config.roundingMode = RoundingModeFloor
}

var FloorBytes = TruncateBytes

func RoundBytes(config *bytesConfig) {
	config.roundingMode = RoundingModeRound
}

func CeilBytes(config *bytesConfig) {
	config.roundingMode = RoundingModeCeil
}

func BytesPrecision(precision int) BytesOption {
	return func(config *bytesConfig) {
		config.precision = precision
	}
}

func TraillingZero(config *bytesConfig) {
	config.traillingZero = true
}

func ValueUnitSeparator(sep string) BytesOption {
	return func(config *bytesConfig) {
		config.valueUnitSeparator = sep
	}
}

func Bytes(numBytes int64, opts ...BytesOption) string {
	config := bytesConfig{
		precision:          1,
		roundingMode:       RoundingModeRound,
		valueUnitSeparator: "",
	}
	for _, opt := range opts {
		opt(&config)
	}

	humanizeFloat := func(value float64) string {
		signal := ""
		if value < 0 {
			signal = "-"
			value = -value
		}
		switch config.roundingMode {
		case RoundingModeCeil:
			value = math.Ceil(value)
		case RoundingModeFloor:
			value = math.Floor(value)
		}
		result := strconv.FormatFloat(value, 'f', config.precision, 64)

		if config.traillingZero {
			return signal + result
		}
		resultBytes := []byte(result)
		for {
			lastChar := resultBytes[len(resultBytes)-1]
			if lastChar != '0' {
				if lastChar == '.' {
					resultBytes = resultBytes[:len(resultBytes)-1]
				}
				break
			}
			resultBytes = resultBytes[:len(resultBytes)-1]
		}
		return signal + string(resultBytes)
	}

	numBytesF := float64(numBytes)
	absNumBytesF := numBytesF
	if absNumBytesF < 0 {
		absNumBytesF = -absNumBytesF
	}
	if absNumBytesF < kib {
		return strconv.FormatInt(numBytes, 10) + config.valueUnitSeparator + "b"
	} else if absNumBytesF < mib {
		return humanizeFloat(numBytesF/kib) + config.valueUnitSeparator + "KiB"
	} else if absNumBytesF < gib {
		return humanizeFloat(numBytesF/mib) + config.valueUnitSeparator + "MiB"
	} else if absNumBytesF < tib {
		return humanizeFloat(numBytesF/gib) + config.valueUnitSeparator + "GiB"
	} else if absNumBytesF < pib {
		return humanizeFloat(numBytesF/tib) + config.valueUnitSeparator + "TiB"
	} else {
		return humanizeFloat(numBytesF/pib) + config.valueUnitSeparator + "PiB"
	}
}
