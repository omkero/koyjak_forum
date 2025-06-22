package functions

import (
	"fmt"
)

func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + " ..."
}

func TruncateFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}

	return s[:1]
}

func CalculateCount(count int) string {

	switch {
	case count >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(count)/1_000_000_000)
	case count >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(count)/1_000_000)
	case count >= 1_000:
		return fmt.Sprintf("%.1fK", float64(count)/1_000)
	case count < 10:
		return fmt.Sprintf("%02d", count)
	default:
		return fmt.Sprintf("%d", count)
	}
}
