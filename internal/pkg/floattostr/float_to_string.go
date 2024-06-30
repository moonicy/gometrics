package floattostr

import "strconv"

func FloatToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', -1, 64)
}
