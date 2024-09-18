// Package floattostr предоставляет функцию для преобразования числа с плавающей точкой в строку.
package floattostr

import "strconv"

// FloatToString преобразует число типа float64 в строку.
func FloatToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', -1, 64)
}
