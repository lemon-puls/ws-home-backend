package cosutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractKeyFromUrl(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "标准 URL",
			url:      "https://www.example.com/exampleobject/1745647348066-761.jpg",
			expected: "exampleobject/1745647348066-761.jpg",
		},
		{
			name:     "带查询参数的 URL",
			url:      "https://www.example.com/exampleobject/1745647348066-761.jpg?q-sign-algorithm=sha1&q-ak=AKIDc6MDsKXWGm38z432-7823gGhv9D4jANM7e094m",
			expected: "exampleobject/1745647348066-761.jpg",
		},
		{
			name:     "空 URL",
			url:      "",
			expected: "",
		},
		{
			name:     "只有域名的 URL",
			url:      "https://www.example.com",
			expected: "",
		},
		{
			name:     "只有 key",
			url:      "ws-home/ablum/1/1749908604372_0sH9vSF33lZJ3c486bda6e7ab3d0b2187ea697c4302c.jpeg",
			expected: "ws-home/ablum/1/1749908604372_0sH9vSF33lZJ3c486bda6e7ab3d0b2187ea697c4302c.jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractKeyFromUrl(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}
