package ripencc

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ParsedLine struct {
	StartIP net.IP
	Size    int
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

	size, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, fmt.Errorf("invalid size: %s", parts[4])
	}

	fmt.Printf("ParsedLine: code=%s, ip=%s, size=%d\n", parts[1], startIP.String(), size)

	return &ParsedLine{
		StartIP: startIP,
		Size:    size,
	}, nil
}
