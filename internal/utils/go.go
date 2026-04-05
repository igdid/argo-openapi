// Written by ChatGPT
package utils

import (
	"runtime"
	"strings"
)

// CurrentGoVersion returns a version suitable for the `go` directive in go.mod,
// e.g. "1.24" from runtime.Version() == "go1.24.1".
func CurrentGoVersion() string {
	v := runtime.Version()

	// Strip "go" prefix: "go1.24.1" -> "1.24.1"
	v = strings.TrimPrefix(v, "go")

	// Handle devel versions like "devel go1.25-abcdef"
	if strings.HasPrefix(v, "devel ") {
		parts := strings.Fields(v)
		if len(parts) > 1 {
			v = strings.TrimPrefix(parts[1], "go")
		}
	}

	parts := strings.Split(v, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}

	return v
}
