package options

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

// Options holds the parsed command-line options
type Options struct {
	Country       string
	Format        string
	Gateway       string
	ListCountries bool
}

// NewOptions creates a new Options instance and parses command-line flags
func NewOptions() (*Options, error) {
	opts := &Options{}

	flag.StringVar(&opts.Country, "c", "", "Country code (required, e.g., RU, US)")
	flag.StringVar(&opts.Format, "f", "plain", "Output format: 'plain', 'ros' or 'ros-list'")
	flag.StringVar(&opts.Gateway, "g", "", "Gateway IP for routes (required with -f=ros)")
	flag.BoolVar(&opts.ListCountries, "list-countries", false, "Print all supported country codes")

	flag.Usage = func() {
		printUsage()
	}

	flag.Parse()

	if opts.ListCountries {
		printSupportedCountryCodes()
	}

	if err := opts.validateCountryCode(); err != nil {
		return nil, err
	}

	if opts.Format == "ros" {
		if err := opts.validateGatewayIP(); err != nil {
			return nil, err
		}
	}

	err := setFormat(opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

// validateCountryCode validates the provided country code
func (opts *Options) validateCountryCode() error {
	if opts.Country == "" {
		return fmt.Errorf("country code is required. Use -c to specify (e.g., -c RU)")
	}
	if _, ok := validCountries[opts.Country]; !ok {
		return fmt.Errorf("error: '%s' is not a supported country code.\nRun '%s -list-countries' to see all options", opts.Country, os.Args[0])
	}

	return nil
}

// validateGatewayIP validates the provided gateway IP address
func (opts *Options) validateGatewayIP() error {
	if opts.Gateway == "" {
		return fmt.Errorf("gateway is required when using -f=ros. Use -g to specify the gateway")
	}
	if net.ParseIP(opts.Gateway) == nil {
		return fmt.Errorf("invalid gateway IP address: %s", opts.Gateway)
	}

	return nil
}

// printUsage prints the usage information for the command-line flags
func printUsage() {
	lines := []string{
		"Usage:",
		fmt.Sprintf("  %s -c <country> [-f <format>] [-g <gateway>]", os.Args[0]),
		"",
		"Flags:",
		flagString(),
		"Examples:",
		fmt.Sprintf("  %s -c US", os.Args[0]),
		"\tOutput routes for the US in plain format",
		fmt.Sprintf("  %s -c US -f ros-list", os.Args[0]),
		"\tOutput routes for the US in RouterOS format to add in address-list",
		fmt.Sprintf("  %s -c RU -f ros -g 192.168.1.1", os.Args[0]),
		"\tOutput routes for Russia in RouterOS format with gateway 192.168.1.1",
		fmt.Sprintf("  %s -list-countries", os.Args[0]),
		"\tList all supported country codes",
		"",
	}
	fmt.Fprintln(os.Stderr, strings.Join(lines, "\n"))
}

// printSupportedCountryCodes prints all supported country codes
func printSupportedCountryCodes() {
	fmt.Println("Supported country codes:")
	fmt.Println(strings.Join(getCountryCodes(), ", "))
	fmt.Printf("\nFor details see: %s\n", countryCodeReference)
	os.Exit(0)
}

// flagString helper function for flag output
func flagString() string {
	var buf strings.Builder
	flag.CommandLine.SetOutput(&buf)
	flag.PrintDefaults()
	return buf.String()
}

func setFormat(opts *Options) error {
	switch opts.Format {
	case "plain":
		opts.Format = "%s/%d\n"
	case "ros":
		opts.Format = fmt.Sprintf("/ip route add dst-address=%%s/%%d gateway=%s comment=\"%s-routes\"\n", opts.Gateway, strings.ToLower(opts.Country))
	case "ros-list":
		opts.Format = fmt.Sprintf("/ip firewall address-list add list=%s-cidrs address=%%s/%%d comment=\"georoute\"\n", strings.ToLower(opts.Country))
	default:
		return fmt.Errorf("invalid format specified: %s. Use 'plain', 'ros' or 'ros-list'", opts.Format)
	}
	return nil
}
