package output

import "strings"

func sanitizePathForInjection(path string) string {
	return strings.ReplaceAll(path, "\\", "\\\\")
}
