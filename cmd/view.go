package cmd

import (
	"fmt"

	"github.com/fatih/color"
)

// output render string
var output = func(datas string) {
	fmt.Print(datas)
}

// perror outputs a red message error from an error
var perror = func(err error) {
	color.Red(err.Error())
}

// success outputs a green successful message
func success(message string) string {
	return color.GreenString(message)
}

// info outputs a blue info message
var info = func(message string) {
	color.Cyan(message)
}
