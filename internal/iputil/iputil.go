package iputil

import (
	"fmt"
	"math/bits"
	"net"
)

// Uint32ToIP converts a uint32 to a net.IP (IPv4)
func Uint32ToIP(uintValue uint32) net.IP {
	return net.IPv4(byte(uintValue>>24), byte((uintValue>>16)&0xFF), byte((uintValue>>8)&0xFF), byte(uintValue&0xFF))
}

// IpToUint32 converts a net.IP (IPv4) to a uint32
func IpToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// CalcPrefix calculates the largest possible CIDR prefix length for a given range
func CalcPrefix(start, end uint32) int {
	if start == end {
		return 32
	}
	size := end - start + 1
	power := uint32(bits.TrailingZeros32(size / 256)) // log2 of largest power of 2 in size
	prefix := 24 - int(power)
	fmt.Printf("CalcPrefix: start=%d, end=%d, size=%d, prefix=%d\n", start, end, size, prefix)
	return prefix
}

// Log2Ceiling calculates the ceiling of log base 2 of n
func Log2Ceiling(n uint32) int {
	if n == 0 {
		return 0
	}
	return bits.Len32(n - 1)
}
