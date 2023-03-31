package main

import (
	"testing"
)

func TestValidConfig(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		flagfields FlagFields
		expected   bool
	}{
		{
			name: "All fields filled",
			flagfields: FlagFields{
				robots_file_path: "value1",
				lang:             "value2",
				api_url:          "value3",
				sitemap_index:    "abxx",
				scroll_timeout:   "1000",
				scroll_size:      2,
			},
			expected: true,
		},
		{
			name: "Empty field",
			flagfields: FlagFields{
				robots_file_path: "value1",
				lang:             "value2",
				api_url:          "value3",
				sitemap_index:    "",
				scroll_timeout:   "",
				scroll_size:      2,
			},
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := validConfig(&tc.flagfields)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
