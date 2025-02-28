package options

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Options holds the parsed command-line options
type Options struct {
	Country string
	Format  string
}

// ParseFlags parses command-line flags and returns options
func ParseFlags() *Options {
	var (
		country, format, gateway string
	)

	flag.StringVar(&country, "c", "", "Country code (required, e.g., RU, US)")
	flag.StringVar(&format, "f", "plain", "Output format: 'plain' or 'ros'")
	flag.StringVar(&gateway, "g", "", "Gateway IP for routes (required with -f=ros)")
	listCountries := flag.Bool("list-countries", false, "Print all supported country codes")

	flag.Usage = func() {
		lines := []string{
			"Usage:",
			fmt.Sprintf("  %s -c <country> [-f <format>] [-g <gateway>]", os.Args[0]),
			"",
			"Flags:",
			flagString(),
			"Examples:",
			fmt.Sprintf("  %s -c US", os.Args[0]),
			"\tOutput routes for the US in plain format",
			fmt.Sprintf("  %s -c RU -f ros -g 192.168.1.1", os.Args[0]),
			"\tOutput routes for Russia in RouterOS format with gateway 192.168.1.1",
			fmt.Sprintf("  %s -list-countries", os.Args[0]),
			"\tList all supported country codes",
			"",
		}
		fmt.Fprintln(os.Stderr, strings.Join(lines, "\n"))
	}

	flag.Parse()

	if *listCountries {
		fmt.Println("Supported country codes:")
		fmt.Println(strings.Join(getCountryCodes(), ", "))
		fmt.Printf("\nFor details see: %s\n", countryCodeReference)
		os.Exit(0)
	}

	// Validate country code (required)
	if country == "" {
		log.Fatalf("Country code is required. Use -c to specify (e.g., -c RU).\n")
	}
	if _, ok := validCountries[country]; !ok {
		log.Fatalf("Error: '%s' is not a supported country code.\nRun '%s -list-countries' to see all options.\n", country, os.Args[0])
	}

	// Prepare options
	opts := &Options{
		Country: country,
	}

	// Set format string based on -f and -g
	switch format {
	case "plain":
		opts.Format = "%s/%d\n"
	case "ros":
		if gateway == "" {
			log.Fatalf("Gateway is required when using -f=ros. Use -g to specify the gateway.")
		}
		if net.ParseIP(gateway) == nil {
			log.Fatalf("Invalid gateway IP address: %s", gateway)
		}
		opts.Format = fmt.Sprintf("/ip route add dst-address=%%s/%%d gateway=%s comment=\"%s-routes\"\n", gateway, strings.ToLower(country))
	default:
		log.Fatalf("Invalid format specified: %s. Use 'plain' or 'ros'.", format)
	}

	return opts
}

// flagString helper function for flag output
func flagString() string {
	var buf strings.Builder
	flag.CommandLine.SetOutput(&buf)
	flag.PrintDefaults()
	return buf.String()
}
