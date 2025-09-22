// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package str

import (
	"regexp"
	"strings"
)

// Convert a string to kebab-case.
func KebabCase(str string) string {
	// Replace non-alphanumeric characters with hyphens
	reAlnumdash := regexp.MustCompile(`[^a-zA-Z0-9-]`)
	str = reAlnumdash.ReplaceAllString(str, "-")
	// Insert hyphen between lowercase and uppercase letters
	reHyphenBetween := regexp.MustCompile(`([a-z])([A-Z])`)
	str = reHyphenBetween.ReplaceAllString(str, "$1-$2")
	// Replace whitespace with hyphens
	reReplaceSpaces := regexp.MustCompile(`[\s]+`)
	str = reReplaceSpaces.ReplaceAllString(str, "-")
	// Replace multiple hyphens with a single hyphen
	reStripHyphens := regexp.MustCompile(`[-]+`)
	str = reStripHyphens.ReplaceAllString(str, "-")
	// Convert to lowercase
	str = strings.ToLower(str)
	// Trim leading hyphens
	str = strings.TrimPrefix(str, "-")
	// Trim trailing hyphens
	str = strings.TrimSuffix(str, "-")

	return str
}
