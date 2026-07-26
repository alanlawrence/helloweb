package main

import (
    "testing"
)

// Provide a test function so this package won't cause an error when
// go test ./... is run from the build root directory.
func TestPrime(t *testing.T) {
    if TestIsPrime() != true {
        t.Errorf("Expected true but got /%v/", false)
    }
}
