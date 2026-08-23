package ripencc

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ParsedLine struct {
	StartIP net.IP
	Size    uint64
}

func ParseLine(code, line string) (*ParsedLine, error) {
	parts := strings.SplitN(strings.TrimSpace(line), "|", 7)
	if len(parts) != 7 || parts[1] != code || parts[2] != "ipv4" {
		return nil, fmt.Errorf("invalid line format")
	}

	startIP := net.ParseIP(parts[3])
	if startIP == nil {
		return nil, fmt.Errorf("invalid IP address: %s", parts[3])
	}

	size, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil || size == 0 {
		return nil, fmt.Errorf("invalid size: %s", parts[4])
	}

	return &ParsedLine{
		StartIP: startIP,
		Size:    size,
	}, nil
}
