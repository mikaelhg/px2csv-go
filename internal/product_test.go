package internal_test

import (
	"testing"

	"github.com/mikaelhg/gpcaxis/internal"
	"github.com/stretchr/testify/assert"
)

func TestCartesianProduct(t *testing.T) {
	data := [][]string{
		{"a", "b", "c"},
		{"1", "2", "3"},
	}
	c := internal.NewCartesianProduct(data)
	result := c.All()
	expected := [][]string{
		{"a", "1"}, {"a", "2"}, {"a", "3"},
		{"b", "1"}, {"b", "2"}, {"b", "3"},
		{"c", "1"}, {"c", "2"}, {"c", "3"},
	}
	assert.Equal(t, expected, result, "should be the same")
}
