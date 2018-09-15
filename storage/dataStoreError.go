package storage

import (
	"fmt"
)

type DataStoreError struct {
	reason string
}

func (error *DataStoreError) Error() string {
	return fmt.Sprintf("DataStore failure. Reason: %s", error.reason)
}