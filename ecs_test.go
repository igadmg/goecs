package ecs

import (
	"testing"
)

func TestId_String(t *testing.T) {
	tests := []struct {
		id       Id
		expected string
	}{
		{MakeId(1, 1), "Id(1:1:false:false)"},
		{MakeId(0, 0), "Id(0:0:false:false)"},
		{MakeId(12345, 678), "Id(12345:678:false:false)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.id.String(); got != tt.expected {
				t.Errorf("Id.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
