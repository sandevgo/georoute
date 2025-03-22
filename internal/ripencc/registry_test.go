package ripencc

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestRegistry_Process(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		country  string
		expected string
	}{
		{
			name: "NoSuitableBlocks",
			input: `ripencc|RU|ipv6|2001:db8::|128|20120522|allocated
ripencc|GB|ipv4|192.168.1.0|256|20120522|assigned
ripencc|AZ|ipv4|172.16.0.0|2048|20120522|reserved`,
			country:  "RU",
			expected: ``,
		},
		{
			name:    "Single24Block",
			input:   `ripencc|RU|ipv4|192.168.1.0|256|20120522|allocated`,
			country: "RU",
			expected: `192.168.1.0/24
`,
		},
		{
			name:    "Single25Block",
			input:   `ripencc|RU|ipv4|192.168.1.0|128|20120522|allocated`,
			country: "RU",
			expected: `192.168.1.0/25
`,
		},
		{
			name: "NonContiguous24Blocks",
			input: `ripencc|RU|ipv4|192.168.1.0|256|20120522|allocated
ripencc|RU|ipv4|192.168.3.0|256|20120522|allocated
ripencc|RU|ipv4|192.168.5.0|256|20120522|allocated`,
			country: "RU",
			expected: `192.168.1.0/24
192.168.3.0/24
192.168.5.0/24
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			reader := strings.NewReader(tt.input)
			registry := NewRegistry(tt.country, "%s/%d\n")
			registry.body = ioutil.NopCloser(reader)
			registry.writer = &output

			if err := registry.Process(); err != nil {
				t.Errorf("Registry.Process() error = %v", err)
			}

			got := output.String()
			if got != tt.expected {
				t.Errorf("Registry.Process() = \n%q\nexpected \n%q", got, tt.expected)
			}
		})
	}
}
