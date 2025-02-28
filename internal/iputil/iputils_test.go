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

func TestCalcPrefix(t *testing.T) {
	tests := []struct {
		name     string
		start    uint32
		end      uint32
		expected int
	}{
		{"Range of 256 IPs", 0xC0A80101, 0xC0A80200, 24},
		{"Range of 512 IPs", 0xC0A80101, 0xC0A80300, 23},
		{"Range of 1024 IPs", 0xC0A80101, 0xC0A80500, 22},
		{"Range of 2048 IPs", 0xC0A80101, 0xC0A80900, 21},
		{"Range of 4096 IPs", 0xC0A80101, 0xC0A81100, 20},
		{"Range of 8192 IPs", 0xC0A80101, 0xC0A82100, 19},
		{"Range of 16384 IPs", 0xC0A80101, 0xC0A84100, 18},
		{"Range of 32768 IPs", 0xC0A80101, 0xC0A88100, 17},
		{"Range of 65536 IPs", 0xC0A80101, 0xC0A90100, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CalcPrefix(tt.start, tt.end)
			if actual != tt.expected {
				t.Errorf("CalcPrefix(%v, %v) = %v; expected %v", tt.start, tt.end, actual, tt.expected)
			}
		})
	}
}
