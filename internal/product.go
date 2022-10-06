package internal

type CartesianProduct struct {
	length   int
	counters []int
	lengths  []int
	lists    [][]string
}

func NewCartesianProduct(input [][]string) CartesianProduct {
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

func (c *CartesianProduct) Reset() {
	c.counters = make([]int, c.length)
}

func (c *CartesianProduct) Next() ([]string, bool) {
	ret := make([]string, c.length)
	for i := 0; i < c.length; i++ {
		ret[i] = c.lists[i][c.counters[i]]
	}
	counters, stop := c.step()
	c.counters = counters
	return ret, stop
}

func (c *CartesianProduct) step() ([]int, bool) {
	ret := make([]int, c.length)
	copy(ret, c.counters)
	for i := 0; i < c.length; i++ {
		if ret[i] < c.lengths[i]-1 {
			ret[i] += 1
			return ret, false
		} else {
			ret[i] = 0
		}
	}
	return nil, true
}

func (c *CartesianProduct) All() [][]string {
	ret := make([][]string, 0)
	for {
		o, stop := c.Next()
		ret = append(ret, o)
		if stop {
			break
		}
	}
	return ret
}
