package utils

import (
	"crypto/sha256"
	"fmt"
)

func Sha256FromStr(input string) string {
	h := sha256.New()
	h.Write([]byte(input))

	return fmt.Sprintf("%x", h.Sum(nil))
}
