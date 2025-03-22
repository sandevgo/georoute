package ripencc

import (
	"bufio"
	"fmt"
	"github.com/sandevgo/georoute/internal/iputil"
	"io"
	"net/http"
	"os"
)

const DelegatedLatest = "https://ftp.ripe.net/pub/stats/ripencc/delegated-ripencc-latest"

type Registry struct {
	country string
	format  string
	body    io.ReadCloser
	writer  io.Writer
}

func NewRegistry(country, format string) *Registry {
	return &Registry{
		country: country,
		format:  format,
		writer:  os.Stdout,
	}
}

func (r *Registry) GetDelegated() error {
	resp, err := http.Get(DelegatedLatest)
	if err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return fmt.Errorf("bad HTTP status: %s", resp.Status)
	}
	r.body = resp.Body
	return nil
}

func (r *Registry) Close() error {
	if r.body != nil {
		return r.body.Close()
	}
	return nil
}

// Process processes the Ripe Delegated data from the reader and writes the CIDR blocks to the writer
func (r *Registry) Process() error {
	scanner := bufio.NewScanner(r.body)
	for scanner.Scan() {
		parsed, err := ParseLine(r.country, scanner.Text())
		if err != nil {
			continue
		}
		
		// Calculate the prefix directly from the size
		// For a block of size N, the prefix is 32 - log2(N)
		var prefix int
		switch {
		case parsed.Size == 1:
			prefix = 32
		case parsed.Size <= 256:
			prefix = 24
		case parsed.Size <= 65536:
			prefix = 16
		case parsed.Size <= 16777216:
			prefix = 8
		default:
			prefix = 0
		}
		
		// Dump each block immediately without aggregation
		fmt.Fprintf(r.writer, r.format, parsed.StartIP.String(), prefix)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading HTTP response: %v", err)
	}

	return nil
}

