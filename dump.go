package main

import (
	"fmt"
)

func dumpCmd() command {
	return command{fn: func([]string) error {
		return fmt.Errorf("trunks dump has been deprecated and succeeded by the trunks encode command")
	}}
}
