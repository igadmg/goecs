package ecs

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestId_String(t *testing.T) {
	tests := []struct {
		id       Id
		expected string
	}{
		{MakeId(1, 1), "Id(1:1:false:false)"},
		{MakeId(0, 0), "Id(0:0:false:false)"},
		{MakeId(12345, 3), "Id(12345:3:false:false)"},
		{MakeId(1, 1).Allocate(), "Id(1:1:true:false)"},
		{MakeId(12345, 3).Store(), "Id(12345:3:false:true)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.id.String(); got != tt.expected {
				t.Errorf("Id.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_RepackList(t *testing.T) {
	free_ids := []int{10, 3, 7, 4, 16}

	slices.Sort(free_ids)
	i, ok := slices.BinarySearch(free_ids, 6)
	assert.Equal(t, 2, i)
	assert.False(t, ok)
}
