// Package retry provides HTTP retry helpers.
package retry

import "fmt"

// PrintError prints debug info for an error.
func PrintError(err error) {
	fmt.Printf("DEBUG: error type: %T, value: %v\n", err, err)
}
