package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortMapValues(t *testing.T) {
	type args struct {
		input map[string]interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected map[string]interface{}
	}{
		{
			name: "sort string slice",
			args: args{
				input: map[string]interface{}{
					"key1": []string{"c", "a", "b"},
				},
			},
			expected: map[string]interface{}{
				"key1": []string{"a", "b", "c"},
			},
		},
		{
			name: "sort int slice",
			args: args{
				input: map[string]interface{}{
					"key1": []int{3, 1, 2},
				},
			},
			expected: map[string]interface{}{
				"key1": []int{1, 2, 3},
			},
		},
		{
			name: "sort map in slice",
			args: args{
				input: map[string]interface{}{
					"key1": []map[string]interface{}{
						{"sort_key": "c"},
						{"sort_key": "a"},
						{"sort_key": "b"},
					},
				},
			},
			expected: map[string]interface{}{
				"key1": []map[string]interface{}{
					{"sort_key": "a"},
					{"sort_key": "b"},
					{"sort_key": "c"},
				},
			},
		},
		{
			name: "nested map",
			args: args{
				input: map[string]interface{}{
					"key1": map[string]interface{}{
						"key11": "a",
						"key12": []string{"c", "b", "a"},
						"key13": map[string]interface{}{
							"key131": []int{2, 3, 1},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"key1": map[string]interface{}{
					"key11": "a",
					"key12": []string{"a", "b", "c"},
					"key13": map[string]interface{}{
						"key131": []int{1, 2, 3},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.expected, SortMapValues(tt.args.input))
		})
	}
}
