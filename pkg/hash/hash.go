// Package hash предоставляет функцию для вычисления SHA256-хэша с использованием заданного ключа.
package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// CalcHash вычисляет SHA256-хэш от переданных данных и ключа.
// Возвращает хэш в виде шестнадцатеречной строки.
func CalcHash(body []byte, key string) string {
	sha := sha256.New()
	sha.Write([]byte(key))
	sha.Write(body)
	shaSum := sha.Sum(nil)
	return hex.EncodeToString(shaSum)
}
