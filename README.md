# GeoRoute

`GeoRoute` is a lightweight command-line tool that generates network routes for a specified country, outputting them in either plain text (CIDR notation) or RouterOS-compatible format. It supports a comprehensive list of ISO 3166-1 alpha-2 country codes and is designed for simplicity and ease of use.

## Features
- Generate routes for any supported country using two-letter country codes (e.g., `US`, `RU`).
- Output formats:
    - `plain`: Simple CIDR notation (e.g., `192.168.1.0/24`).
    - `ros`: RouterOS commands with a specified gateway (e.g., `/ip route add dst-address=192.168.1.0/24 gateway=192.168.1.1`).
- List all supported country codes with a single flag.
- Built-in validation for country codes and gateway IPs.

## Installation

### Prerequisites
- [Go](https://golang.org/dl/) 1.24 or later.

### Build from Source
1. Clone the repository:
   ```bash
   git clone https://github.com/sandevgo/georoute.git
   cd georoute
   ```
2. Build the binary:
   ```bash
   make
   ```
3. (Optional) Move the binary to a directory in your PATH:
   ```bash
   sudo mv georoute /usr/local/bin/
   ```

## Usage
Run the program with the required `-c` flag to specify a country code. Use `-f` to set the output format and `-g` for the gateway IP (required with `-f ros`).

```
georoute -c <country> [-f <format>] [-g <gateway>]
```

### Flags
- `-c <country>`: Country code (required, e.g., `RU`, `US`).
- `-f <format>`: Output format: `plain` (default) or `ros`.
- `-g <gateway>`: Gateway IP for routes (required with `-f ros`).
- `-list-countries`: Print all supported country codes.
- `-h`: Show help with usage and examples.

### Examples
1. Generate routes for the United States in plain format:
   ```bash
   georoute -c US
   ```
   Output:
   ```
   192.0.2.0/24
   198.51.100.0/24
   ...
   ```

2. Generate routes for Russia in RouterOS format with a gateway:
   ```bash
   georoute -c RU -f ros -g 192.168.1.1
   ```
   Output:
   ```
   /ip route add dst-address=93.158.0.0/18 gateway=192.168.1.1 comment="georoute"
   /ip route add dst-address=94.100.176.0/20 gateway=192.168.1.1 comment="georoute"
   ...
   ```

3. List all supported country codes:
   ```bash
   georoute -list-countries
   ```
   Output:
   ```
   Supported country codes:
   AD, AE, AF, AG, AI, AL, AM, AO, AR, AT, AU, AX, AZ, BA, BD, BE, BG, BH, BI, BM, BO, BQ, BR, BY, BZ, CA, ...

   For details see: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
   ```

4. View help:
   ```bash
   georoute -h
   ```
   Output:
   ```
   Usage:
     georoute -c <country> [-f <format>] [-g <gateway>]

   Flags:
     -c string
           Country code (required, e.g., RU, US)
     -f string
           Output format: 'plain' or 'ros' (default "plain")
     -g string
           Gateway IP for routes (required with -f=ros)
     -list-countries
           Print all supported country codes

   Examples:
     georoute -c US
           Output routes for the US in plain format
     georoute -c RU -f ros -g 192.168.1.1
           Output routes for Russia in RouterOS format with gateway 192.168.1.1
     georoute -list-countries
           List all supported country codes
   ```

## Notes
- Country codes follow the [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) standard.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
