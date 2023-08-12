package cmd

import (
	"github.com/fatih/color"
)

// info outputs a blue info message
func info(message string) string {
	return color.CyanString(message)
}

// success outputs a green successful message
func success(message string) string {
	return color.GreenString(message)
}
