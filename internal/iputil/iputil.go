package iputil

import (
	"fmt"
	"math/bits"
	"net"
)

// RangeToCIDRs converts an IPv4 address range to its minimal CIDR cover.
func RangeToCIDRs(startIP net.IP, size uint64) ([]net.IPNet, error) {
	startIP = startIP.To4()
	if startIP == nil {
		return nil, fmt.Errorf("start IP is not IPv4")
	}
	if size == 0 {
		return nil, fmt.Errorf("range size must be positive")
	}

	start := IpToUint32(startIP)
	remaining := size
	if remaining > (uint64(1)<<32)-uint64(start) {
		return nil, fmt.Errorf("IPv4 range exceeds 255.255.255.255")
	}

	blocks := make([]net.IPNet, 0, 32)
	current := start
	for remaining > 0 {
		maxByAlignment := uint64(1) << bits.TrailingZeros32(current)
		maxByRemaining := uint64(1) << (bits.Len64(remaining) - 1)
		blockSize := min(maxByAlignment, maxByRemaining)
		prefix := 32 - bits.TrailingZeros64(blockSize)

		blocks = append(blocks, net.IPNet{
			IP:   Uint32ToIP(current),
			Mask: net.CIDRMask(prefix, 32),
		})

		remaining -= blockSize
		if remaining > 0 {
			current = uint32(uint64(current) + blockSize)
		}
	}

	return blocks, nil
}

// Uint32ToIP converts a uint32 to a net.IP (IPv4)
func Uint32ToIP(uintValue uint32) net.IP {
	return net.IPv4(byte(uintValue>>24), byte((uintValue>>16)&0xFF), byte((uintValue>>8)&0xFF), byte(uintValue&0xFF))
}

// IpToUint32 converts a net.IP (IPv4) to a uint32
func IpToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
