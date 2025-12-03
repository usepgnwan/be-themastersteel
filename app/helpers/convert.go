package helpers

import "strconv"

func ParseFloat(val string) float64 {
	f, _ := strconv.ParseFloat(val, 64)
	return f
}

func ParseInt(val string) int {
	i, _ := strconv.Atoi(val)
	return i
}
