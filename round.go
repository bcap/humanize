package humanize

type roundingMode int

const (
	RoundingModeFloor roundingMode = -1
	RoundingModeRound roundingMode = 0
	RoundingModeCeil  roundingMode = 1
)
