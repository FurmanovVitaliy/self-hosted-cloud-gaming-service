package hashsum

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func Hashsum(obj ...any) string {
	hash := sha256.New()
	for _, o := range obj {
		hash.Write([]byte(fmt.Sprintf("%v", o)))
	}
	hashSum := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashSum)
}

func CheckSum(first string, second string) bool {
	if first == "" || second == "" {
		return false
	}
	if first == second {
		return true
	}
	return false
}
