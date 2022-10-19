package internal

import "strings"

func MapXtoY[X interface{}, Y interface{}](collection []X, fn func(elem X) Y) []Y {
	var result []Y
	for _, item := range collection {
		result = append(result, fn(item))
	}
	return result
}

func joinStringSlice(ss []string) string {
	return strings.Join(ss, " ")
}
