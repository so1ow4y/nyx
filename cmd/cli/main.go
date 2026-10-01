package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"gitlab.com/ere6u5/nyx/internal/config"
	"gitlab.com/ere6u5/nyx/internal/core"
	"gitlab.com/ere6u5/nyx/internal/output"
	"gitlab.com/ere6u5/nyx/internal/parser"
	"gitlab.com/ere6u5/nyx/pkg/constants"
	"gitlab.com/ere6u5/nyx/pkg/types"
)

func main() {

	configuration := parseFlags()

	parsedPorts, err := parser.ParsePorts(configuration.Ports)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ports: %v\n", err)
		os.Exit(1)
	}

	if configuration.Host == "" {
		fmt.Fprintf(os.Stdout, "Usage: %s --help for check examples\n", constants.NameBin)
		os.Exit(0)
	}

	if _, err := net.LookupHost(configuration.Host); err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid host %s - %v\n", configuration.Host, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Starting %s %s\n", constants.NameApp, constants.VersionApp)
	fmt.Fprintf(os.Stdout, "Scanning %s (%s)\n", configuration.Host, configuration.Protocol)

	results := core.Scan(configuration.Host, parsedPorts, configuration.Protocol, configuration.Timeout)

	if configuration.OutputJSON == "" {
		for _, checkedPort := range results {
			fmt.Fprintf(os.Stdout, "  %d - %s\n", checkedPort.Port, checkedPort.Status)
		}
	} else {
		err := output.SaveToJSONFile(configuration.OutputJSON, results)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Results saved to %s\n", configuration.OutputJSON)
	}

	fmt.Fprintf(os.Stdout, "Finish\n")
}

// TODO doc
func parseFlags() types.Config {
	var configuration types.Config

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "%s %s\n", constants.NameApp, constants.VersionApp)
		fmt.Fprintf(os.Stdout, "A simple port scanner.\n")
		fmt.Fprintf(os.Stdout, "\nUsage:\n")
		fmt.Fprintf(os.Stdout, "    %s -h <host> [-p <ports>] [-tp <protocol>]\n", constants.NameBin)

		fmt.Fprintf(os.Stdout, "\nOptions:\n")
		fmt.Fprintf(os.Stdout, "  -h, --host: Ip address or url to host.\n")
		fmt.Fprintf(os.Stdout, "    Can use ip, ip4 (IPv4-only), ip6 (IPv6-only), url.\n")
		fmt.Fprintf(os.Stdout, "  -p, --ports <port range>: Only scan specified ports.\n")
		fmt.Fprintf(os.Stdout, "  -tp, --type-protocol: Type of internet protocol.\n")
		fmt.Fprintf(os.Stdout, "    All types: tcp, tcp4 (IPv4-only), tcp6 (IPv6-only), udp, udp4 (IPv4-only), udp6 (IPv6-only).\n")
		fmt.Fprintf(os.Stdout, "  -jf, --json-file: Output results in JSON format\n")
		fmt.Fprintf(os.Stdout, "  -to, --timeout: Connection timeout (e.g. 2s, 500ms)\n")

		fmt.Fprintf(os.Stdout, "\nExamples:\n")
		fmt.Fprintf(os.Stdout, "  %s -h example.com -p 443,54,70-80\n", constants.NameBin)
		fmt.Fprintf(os.Stdout, "  %s -h=example.com --ports 443,54,70-80\n", constants.NameBin)
		fmt.Fprintf(os.Stdout, "  %s --host=example.com -p=443,54,70-80 -jf out.json\n", constants.NameBin)
		fmt.Fprintf(os.Stdout, "  %s --host=example.com --ports=443,54,70-80 --json-file out.json\n", constants.NameBin)
	}

	flag.StringVar(&configuration.Host, "h", "", "Target host (IP or domain)")
	flag.StringVar(&configuration.Host, "host", "", "Target host (IP or domain)")
	flag.StringVar(&configuration.Ports, "p", config.DefaultPorts, "Ports to scan (comma-separated or range)")
	flag.StringVar(&configuration.Ports, "ports", config.DefaultPorts, "Ports to scan")
	flag.StringVar(&configuration.Protocol, "tp", config.DefaultProtocol, "Protocol type (tcp/udp/etc.)")
	flag.StringVar(&configuration.Protocol, "type-protocol", config.DefaultProtocol, "Protocol type")
	flag.StringVar(&configuration.OutputJSON, "jf", "", "Output results in JSON format")
	flag.StringVar(&configuration.OutputJSON, "json-file", "", "Output results in JSON format")
	flag.DurationVar(&configuration.Timeout, "to", config.DefaultTimeout, "Connection timeout (e.g. 2s, 500ms)")
	flag.DurationVar(&configuration.Timeout, "timeout", config.DefaultTimeout, "Connection timeout (e.g. 2s, 500ms)")

	flag.Parse()

	if configuration.Timeout < 100*time.Millisecond {
		fmt.Println("Warning: timeout too low, setting to minimum 100ms")
		configuration.Timeout = 100 * time.Millisecond
	}

	if configuration.Timeout > 10*time.Second {
		fmt.Println("Warning: timeout too high, setting to maximum 10s")
		configuration.Timeout = 10 * time.Second
	}

	return configuration
}
