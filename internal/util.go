package internal

func MapXtoY[X interface{}, Y interface{}](collection []X, fn func(elem X) Y) []Y {
	var result []Y
	for _, item := range collection {
		result = append(result, fn(item))
	}
	return result
}
