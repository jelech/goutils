package goutils

import "fmt"

// Version represents the version of goutils
const Version = "v1.0.0"

// GetVersion returns the current version of goutils
func GetVersion() string {
	return Version
}

// Hello returns a hello message with version info
func Hello() string {
	return fmt.Sprintf("Welcome to GoUtils %s - A collection of useful Go utilities", Version)
}
