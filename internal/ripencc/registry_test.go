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
		{
			name: "ContiguousRanges",
			input: `ripencc|RU|ipv4|5.44.0.0|4096|20120522|allocated
ripencc|GB|ipv4|5.44.16.0|4096|20120522|allocated
ripencc|AZ|ipv4|5.44.32.0|2048|20120522|allocated
ripencc|RU|ipv4|5.44.40.0|512|20120523|allocated
ripencc|RU|ipv4|5.44.42.0|256|20120523|allocated
ripencc|RU|ipv4|5.44.43.0|256|20120523|allocated
ripencc|RU|ipv4|5.44.44.0|256|20120523|allocated
ripencc|RU|ipv4|5.44.45.0|256|20120523|allocated
ripencc|RU|ipv4|5.44.46.0|512|20120523|allocated
ripencc|RU|ipv4|5.44.48.0|4096|20120523|allocated
ripencc|NO|ipv4|5.44.64.0|2048|20120523|allocated
ripencc|NL|ipv4|5.44.72.0|2048|20120523|allocated
ripencc|TR|ipv4|5.44.80.0|4096|20120523|allocated
ripencc|DE|ipv4|5.44.96.0|4096|20120523|allocated`,
			country: "RU",
			expected: `5.44.0.0/20
5.44.40.0/21
`,
		},
		{
			name: "MultipleContiguousBlocks",
			input: `ripencc|RU|ipv4|5.44.0.0|4096|20120522|allocated
ripencc|GB|ipv4|5.44.16.0|4096|20120522|allocated
ripencc|AZ|ipv4|5.44.32.0|2048|20120522|allocated
ripencc|RU|ipv4|5.44.40.0|512|20120523|allocated
ripencc|RU|ipv4|5.44.42.0|256|20120523|allocated
ripencc|RU|ipv4|5.44.43.0|256|20120523|allocated
ripencc|RU|ipv4|5.44.44.0|256|20120523|allocated
ripencc|RU|ipv4|5.44.45.0|256|20120523|allocated
ripencc|RU|ipv4|5.44.46.0|512|20120523|allocated
ripencc|RU|ipv4|5.44.48.0|4096|20120523|allocated
ripencc|NO|ipv4|5.44.64.0|2048|20120523|allocated
ripencc|NL|ipv4|5.44.72.0|2048|20120523|allocated
ripencc|TR|ipv4|5.44.80.0|4096|20120523|allocated
ripencc|DE|ipv4|5.44.96.0|4096|20120523|allocated
ripencc|RU|ipv4|192.168.1.0|256|20120522|allocated
ripencc|RU|ipv4|192.168.2.0|256|20120522|allocated
ripencc|GB|ipv4|192.168.2.0|256|20120522|allocated`,
			country: "RU",
			expected: `5.44.0.0/20
5.44.40.0/21
192.168.1.0/23
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
