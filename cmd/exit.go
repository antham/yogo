package cmd

import (
	"os"
)

// errorExit exit with error status code 1
var errorExit = func() {
	os.Exit(1)
}

// SuccessExit exit with success status code 0
var successExit = func() {
	os.Exit(0)
}
