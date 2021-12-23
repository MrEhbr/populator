package helpers

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func TryHexStringToBytes(s string) ([]byte, error) {
	if !strings.HasPrefix(s, "0x") {
		return nil, fmt.Errorf("not a hex string")
	}
	return hex.DecodeString(s[2:])
}
