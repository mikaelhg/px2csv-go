package internal

type CartesianProduct struct {
	length   int
	counters []int
	lengths  []int
	lists    [][]interface{}
}

func NewCartesianProduct(input [][]interface{}) CartesianProduct {
	length := len(input)
	ret := CartesianProduct{
		length:   length,
		counters: make([]int, length),
		lengths:  make([]int, length),
		lists:    input,
	}
	for i := 0; i < length; i++ {
		ret.lengths[i] = len(input[i])
	}
	return ret
}

func (c *CartesianProduct) Next() []interface{} {
	ret := make([]interface{}, c.length)
	for i := 0; i < c.length; i++ {
		ret[i] = c.lists[i][c.counters[i]]
	}
	c.counters = c.step()
	return ret
}

func (c *CartesianProduct) step() []int {
	ret := make([]int, c.length)
	copy(ret, c.counters)
	for i := 0; i < c.length; i++ {
		if ret[i] == c.lengths[i] {
			ret[i] = 0
		} else {
			ret[i] += 1
			return ret
		}
	}
	return []int{}
}
