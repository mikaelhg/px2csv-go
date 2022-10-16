package internal_test

import (
	"testing"

	"github.com/mikaelhg/gpcaxis/internal"
	"github.com/stretchr/testify/assert"
)

func TestCartesianProductAll(t *testing.T) {
	data := [][]string{
		{"a", "b", "c"},
		{"1", "2", "3"},
	}
	expected := [][]string{
		{"a", "1"}, {"a", "2"}, {"a", "3"},
		{"b", "1"}, {"b", "2"}, {"b", "3"},
		{"c", "1"}, {"c", "2"}, {"c", "3"},
	}
	c := internal.NewCartesianProduct(data)
	result := c.All()
	assert.Equal(t, expected, result, "should be the same")
}

func TestCartesianProductPointer(t *testing.T) {
	data := [][]string{
		{"a", "b", "c"},
		{"1", "2", "3"},
	}
	expected_v := [][]string{
		{"a", "1"}, {"a", "2"}, {"a", "3"},
		{"b", "1"}, {"b", "2"}, {"b", "3"},
		{"c", "1"}, {"c", "2"}, {"c", "3"},
	}
	expected := make([][]*string, len(expected_v))
	for i := 0; i < len(expected_v); i++ {
		expected[i] = make([]*string, len(expected_v[i]))
		for j := 0; j < len(expected_v[i]); j++ {
			expected[i][j] = &(expected_v[i][j])
		}
	}
	c := internal.NewCartesianProduct(data)
	values := make([]*string, len(data))
	for i := 0; i < len(expected); i++ {
		stop := c.NextP(&values)
		assert.True(t, i < len(expected)-1 != stop, "stop on the last step")
		assert.Equal(t, values, expected[i])
	}
}
