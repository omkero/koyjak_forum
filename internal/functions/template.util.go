package functions

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
