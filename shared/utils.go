package shared

import (
	"fmt"
)

// Function to convert currentBlock to hexadecimal
func CurrentBlockToHex(currentBlock int64) string {
	hexString := fmt.Sprintf("0x%x", currentBlock)
	return hexString
}
