package base36

import (
	"strings"

	b3 "github.com/martinlindhe/base36"
)

func EncodeToStringUC(b []byte) string {
	return b3.EncodeBytes(b)
}
func EncodeToString(b []byte) string {
	return strings.ToLower(b3.EncodeBytes(b))
}
func DecodeString(s string) ([]byte, error) {
	return b3.DecodeToBytes(strings.ToUpper(s)), nil
}
