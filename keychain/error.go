package keychain

import "fmt"

type ErrNotFound struct {
	Key string
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("Record for host %s not found", err.Key)
}
