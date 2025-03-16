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
	var rangeStart, rangeEnd uint32

	scanner := bufio.NewScanner(r.body)
	for scanner.Scan() {
		parsed, err := ParseLine(r.country, scanner.Text())
		if err != nil {
			continue
		}

		startUint := iputil.IpToUint32(parsed.StartIP)
		endUint := startUint + uint32(parsed.Size) - 1

		fmt.Printf("Processing block: start=%s, size=%d, prefix=%d\n",
			parsed.StartIP.String(), parsed.Size, iputil.CalcPrefix(startUint, endUint))

		// First range
		if rangeStart == 0 {
			rangeStart = startUint
			rangeEnd = endUint
			continue
		}

		// Check if the current range is contiguous with the previous range
		if startUint == rangeEnd+1 {
			// Merge the ranges
			rangeEnd = endUint
			fmt.Printf("Merging block: new_end=%s\n", iputil.Uint32ToIP(rangeEnd).String())
		} else {
			// Output the merged range
			r.Dump(rangeStart, rangeEnd)
			rangeStart = startUint
			rangeEnd = endUint
			fmt.Printf("New range: start=%s, end=%s\n",
				iputil.Uint32ToIP(rangeStart).String(), iputil.Uint32ToIP(rangeEnd).String())
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading HTTP response: %v", err)
	}

	// Dump the last range
	r.Dump(rangeStart, rangeEnd)

	return nil
}

func (r *Registry) Dump(rangeStart, rangeEnd uint32) {
	fmt.Fprintf(r.writer, r.format, iputil.Uint32ToIP(rangeStart).String(), iputil.CalcPrefix(rangeStart, rangeEnd))
}
