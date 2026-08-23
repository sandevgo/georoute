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

		blocks, err := iputil.RangeToCIDRs(parsed.StartIP, parsed.Size)
		if err != nil {
			continue
		}

		for _, block := range blocks {
			prefix, _ := block.Mask.Size()
			fmt.Fprintf(r.writer, r.format, block.IP.String(), prefix)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading HTTP response: %v", err)
	}

	return nil
}
