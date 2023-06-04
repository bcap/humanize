package humanize

import "strconv"

func formatFloat(value float64, precision int) string {
	buf := strconv.AppendFloat(nil, value, 'f', precision, 64)
	if precision == 0 {
		return string(buf)
	}
	endIdx := len(buf) - 1
	for buf[endIdx] == '0' {
		endIdx--
	}
	if buf[endIdx] == '.' {
		endIdx--
	}
	return string(buf[:endIdx+1])
}
