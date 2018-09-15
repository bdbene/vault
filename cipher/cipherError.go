package cipher

import (
	"fmt"
)

type CipherError struct {
	action string
	reason string
}

func (error *CipherError) Error() string {
	return fmt.Sprintf("%s failed. Reason: %s", error.action, error.reason)
}