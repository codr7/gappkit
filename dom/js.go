package dom

import (
	"strings"
)

func EscapeString(in string) string {
	return strings.ReplaceAll(in, "'", "\\'")
}
