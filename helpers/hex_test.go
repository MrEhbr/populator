package helpers

import (
	"testing"
)

func TestTryHexStringToBytes(t *testing.T) {
	TryHexStringToBytes("0x12323213213123")
}
