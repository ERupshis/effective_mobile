package helpers

import (
	"fmt"
)

func ExecuteWithLogError(callback func() error) {
	if err := callback(); err != nil {
		//TODO: add log.
		fmt.Println("error")
	}
}
