// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package utils

import (
	"fmt"
	"log"
	"strings"
)

// PrintHeader prints a header message
func PrintHeader(header string) {
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))
}

// ConsoleOutput writes to stdout and to a logger.
func ConsoleOutput(message string) {
	log.Println(message)
}

// Contains checks if there is a string already in an existing splice of strings
func Contains(array []string, element string) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}
	return false
}

// GetBytesInTiB converts a value from bytes to tebibytes (TiB)
func GetBytesInTiB(size uint64) uint32 {
	return uint32(size / 1024 / 1024 / 1024 / 1024)
}

// GetTiBInBytes converts a value from tebibytes (TiB) to bytes
func GetTiBInBytes(size uint32) uint64 {
	return uint64(size * 1024 * 1024 * 1024 * 1024)
}
