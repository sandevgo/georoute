package iputil

import (
	"net"
	"testing"
)

func TestUint32ToIP(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected net.IP
	}{
		{"Test 192.168.1.1", 0xC0A80101, net.ParseIP("192.168.1.1")},
		{"Test 10.0.0.0", 0x0A000000, net.ParseIP("10.0.0.0")},
		{"Test 0.0.0.0", 0x00000000, net.ParseIP("0.0.0.0")},
		{"Test 255.255.255.255", 0xFFFFFFFF, net.ParseIP("255.255.255.255")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Uint32ToIP(tt.input)
			if !actual.Equal(tt.expected) {
				t.Errorf("Uint32ToIP(%v) = %v; expected %v", tt.input, actual.String(), tt.expected.String())
			}
		})
	}
}

func TestIpToUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    net.IP
		expected uint32
	}{
		{"Test 192.168.1.1", net.ParseIP("192.168.1.1"), 0xC0A80101},
		{"Test 10.0.0.0", net.ParseIP("10.0.0.0"), 0x0A000000},
		{"Test 0.0.0.0", net.ParseIP("0.0.0.0"), 0x00000000},
		{"Test 255.255.255.255", net.ParseIP("255.255.255.255"), 0xFFFFFFFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IpToUint32(tt.input)
			if actual != tt.expected {
				t.Errorf("IpToUint32(%v) = %v; expected %v", tt.input.String(), actual, tt.expected)
			}
		})
	}
}

func TestRangeToCIDRs(t *testing.T) {
	tests := []struct {
		name     string
		start    string
		size     uint64
		expected []string
	}{
		{
			name:     "aligned power of two",
			start:    "192.168.1.0",
			size:     256,
			expected: []string{"192.168.1.0/24"},
		},
		{
			name:     "misaligned power of two",
			start:    "31.135.244.0",
			size:     2048,
			expected: []string{"31.135.244.0/22", "31.135.248.0/22"},
		},
		{
			name:     "non power of two",
			start:    "62.122.208.0",
			size:     1280,
			expected: []string{"62.122.208.0/22", "62.122.212.0/24"},
		},
		{
			name:     "last IPv4 address",
			start:    "255.255.255.255",
			size:     1,
			expected: []string{"255.255.255.255/32"},
		},
		{
			name:     "all IPv4 addresses",
			start:    "0.0.0.0",
			size:     uint64(1) << 32,
			expected: []string{"0.0.0.0/0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks, err := RangeToCIDRs(net.ParseIP(tt.start), tt.size)
			if err != nil {
				t.Fatalf("RangeToCIDRs() error = %v", err)
			}
			if len(blocks) != len(tt.expected) {
				t.Fatalf("RangeToCIDRs() returned %d blocks; expected %d", len(blocks), len(tt.expected))
			}
			for i, block := range blocks {
				if actual := block.String(); actual != tt.expected[i] {
					t.Errorf("RangeToCIDRs() block %d = %q; expected %q", i, actual, tt.expected[i])
				}
			}
		})
	}
}

func TestRangeToCIDRsRejectsInvalidRanges(t *testing.T) {
	tests := []struct {
		name  string
		start net.IP
		size  uint64
	}{
		{name: "IPv6 start", start: net.ParseIP("2001:db8::"), size: 1},
		{name: "zero size", start: net.ParseIP("192.0.2.0"), size: 0},
		{name: "range overflow", start: net.ParseIP("255.255.255.255"), size: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := RangeToCIDRs(tt.start, tt.size); err == nil {
				t.Fatal("RangeToCIDRs() error = nil")
			}
		})
	}
}
