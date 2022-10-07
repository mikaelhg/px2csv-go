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
	return ret, c.step()
}

func (c *CartesianProduct) step() bool {
	for i := c.length - 1; i >= 0; i-- {
		if c.counters[i] < c.lengths[i]-1 {
			c.counters[i] += 1
			return false
		} else {
			c.counters[i] = 0
		}
	}
	return true
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
