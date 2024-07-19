package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func CalcHash(body []byte, key string) string {
	sha := sha256.New()
	sha.Write([]byte(key))
	sha.Write(body)
	shaSum := sha.Sum(nil)
	return hex.EncodeToString(shaSum)
}
