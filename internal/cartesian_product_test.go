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
	expectedValues := [][]string{
		{"a", "1"}, {"a", "2"}, {"a", "3"},
		{"b", "1"}, {"b", "2"}, {"b", "3"},
		{"c", "1"}, {"c", "2"}, {"c", "3"},
	}
	expectedPointers := make([][]*string, len(expectedValues))
	for i := 0; i < len(expectedValues); i++ {
		expectedPointers[i] = make([]*string, len(expectedValues[i]))
		for j := 0; j < len(expectedValues[i]); j++ {
			expectedPointers[i][j] = &(expectedValues[i][j])
		}
	}
	c := internal.NewCartesianProduct(data)
	values := make([]*string, len(data))
	for i := 0; i < len(expectedPointers); i++ {
		stop := c.NextP(&values)
		assert.True(t, i < len(expectedPointers)-1 != stop, "stop on the last step")
		assert.Equal(t, values, expectedPointers[i])
	}
}
