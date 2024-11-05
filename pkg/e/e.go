// Package e is a simple package for wrapping errors
package e

import "fmt"

// Wrap wraps an error with an additional message string
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapIfErr does wrapping only if error is not nil
func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	return Wrap(msg, err)
}
