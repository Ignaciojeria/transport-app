package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination_IsValid(t *testing.T) {
	tests := []struct {
		name        string
		pagination  Pagination
		expectError bool
	}{
		{
			name: "valid forward pagination with first",
			pagination: Pagination{
				First: ptr(10),
			},
			expectError: false,
		},
		{
			name: "valid forward pagination with first and after",
			pagination: Pagination{
				First: ptr(10),
				After: strPtr("cursor:10"),
			},
			expectError: false,
		},
		{
			name: "valid backward pagination with last",
			pagination: Pagination{
				Last: ptr(10),
			},
			expectError: false,
		},
		{
			name: "valid backward pagination with last and before",
			pagination: Pagination{
				Last:   ptr(10),
				Before: strPtr("cursor:20"),
			},
			expectError: false,
		},
		{
			name: "invalid: both first and last",
			pagination: Pagination{
				First: ptr(10),
				Last:  ptr(10),
			},
			expectError: true,
		},
		{
			name: "invalid: first with before",
			pagination: Pagination{
				First:  ptr(10),
				Before: strPtr("cursor:20"),
			},
			expectError: true,
		},
		{
			name: "invalid: last with after",
			pagination: Pagination{
				Last:  ptr(10),
				After: strPtr("cursor:10"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pagination.IsValid()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPagination_IsForward(t *testing.T) {
	tests := []struct {
		name       string
		pagination Pagination
		expected   bool
	}{
		{
			name: "forward pagination with first",
			pagination: Pagination{
				First: ptr(10),
			},
			expected: true,
		},
		{
			name: "not forward pagination",
			pagination: Pagination{
				Last: ptr(10),
			},
			expected: false,
		},
		{
			name:       "empty pagination",
			pagination: Pagination{},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pagination.IsForward())
		})
	}
}

func TestPagination_IsBackward(t *testing.T) {
	tests := []struct {
		name       string
		pagination Pagination
		expected   bool
	}{
		{
			name: "backward pagination with last",
			pagination: Pagination{
				Last: ptr(10),
			},
			expected: true,
		},
		{
			name: "not backward pagination",
			pagination: Pagination{
				First: ptr(10),
			},
			expected: false,
		},
		{
			name:       "empty pagination",
			pagination: Pagination{},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pagination.IsBackward())
		})
	}
}

func TestPagination_HasAfter(t *testing.T) {
	tests := []struct {
		name       string
		pagination Pagination
		expected   bool
	}{
		{
			name: "has after cursor",
			pagination: Pagination{
				After: strPtr("cursor:10"),
			},
			expected: true,
		},
		{
			name: "empty after cursor",
			pagination: Pagination{
				After: strPtr(""),
			},
			expected: false,
		},
		{
			name:       "no after cursor",
			pagination: Pagination{},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pagination.HasAfter())
		})
	}
}

func TestPagination_HasBefore(t *testing.T) {
	tests := []struct {
		name       string
		pagination Pagination
		expected   bool
	}{
		{
			name: "has before cursor",
			pagination: Pagination{
				Before: strPtr("cursor:20"),
			},
			expected: true,
		},
		{
			name: "empty before cursor",
			pagination: Pagination{
				Before: strPtr(""),
			},
			expected: false,
		},
		{
			name:       "no before cursor",
			pagination: Pagination{},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pagination.HasBefore())
		})
	}
}

// Helper function to create pointer to int
func ptr(i int) *int {
	return &i
}

// Helper function to create pointer to string
func strPtr(s string) *string {
	return &s
}
