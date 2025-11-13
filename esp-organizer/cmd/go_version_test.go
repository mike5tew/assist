package main

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGoVersion(t *testing.T) {
	fmt.Printf("Current Go version: %s\n", runtime.Version())
}
