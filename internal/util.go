package internal

import "strings"

func joinStringSlice(ss []string, _ int) string {
	return strings.Join(ss, " ")
}
