package fieldutils

import "strings"

func NormalizedStringField(str string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(str), ""))
}
