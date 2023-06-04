package humanize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	test := func(value int64, expected string, opts ...BytesOption) {
		assert.Equal(t, expected, Bytes(value, opts...))
		assert.Equal(t, "-"+expected, Bytes(-value, opts...))
	}

	test(1, "1b")
	test(1024, "1KiB")
	test(1024*1024, "1MiB")
	test(1024*1024*1024, "1GiB")
	test(1024*1024*1024*1024, "1TiB")
	test(1024*1024*1024*1024*1024, "1PiB")
	test(1024*1024*1024*1024*1024*1024, "1024PiB")

	test(198, "198b")
	test(1023, "1023b")

	test(1024*1024-1, "1024KiB")
	test(1024*1024-1, "1024KiB", BytesPrecision(1))
	test(1024*1024-1, "1024.0KiB", BytesPrecision(1), TraillingZero)
	test(1024*1024-1, "1024KiB", BytesPrecision(2))
	test(1024*1024-1, "1024.00KiB", BytesPrecision(2), TraillingZero)
	test(1024*1024-1, "1023.999KiB", BytesPrecision(3))
	test(1024*1024-1, "1023.999KiB", BytesPrecision(3), TraillingZero)
	test(1024*1024-1, "1023.999KiB", BytesPrecision(4))
	test(1024*1024-1, "1023.9990KiB", BytesPrecision(4), TraillingZero)

	test(1024+1024/3, "1KiB", BytesPrecision(0))
	test(1024+1024/3, "1.3KiB", BytesPrecision(1))
	test(1024+1024/3, "1.33KiB", BytesPrecision(2))
	test(1024+1024/3, "1.333KiB", BytesPrecision(3))

	test(1024+1024/3*2, "1.7KiB", RoundBytes)
	test(1024+1024/3*2, "2KiB", CeilBytes)
	test(1024+1024/3*2, "1KiB", FloorBytes)
	test(1024+1024/3*2, "1KiB", TruncateBytes)
}
