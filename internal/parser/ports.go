package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/ere6u5/nyx/pkg/constants"
)

// ParsePorts: TODO doc
func ParsePorts(inputPorts string) ([]int, error) {

	// check for invalids characters
	if matched, _ := regexp.MatchString(`[^0-9,\-]`, inputPorts); matched {
		return nil, fmt.Errorf("invalid characters in port specification")
	}

	var ports []int
	portsRanges := strings.Split(inputPorts, ",")

	for _, pr := range portsRanges {

		if pr == "" {
			continue
		}

		if strings.Contains(pr, "-") {
			rangePorts, err := ParseRangePorts(pr)
			if err != nil {
				return nil, err
			}

			ports = append(ports, rangePorts...)
		} else {
			port, err := ParseSinglePort(pr)
			if err != nil {
				return nil, err
			}

			ports = append(ports, port)
		}

	}

	return DeduplicatePorts(ports), nil
}

// TODO doc
func ParseRangePorts(pr string) ([]int, error) {
	// handle port range
	parts := strings.Split(pr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invaliv port range format: %s", pr)
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid start port in range %s: %v", pr, err)
	}
	if start < 0 || start > constants.MaxPort {
		return nil, fmt.Errorf("start port out of range (0-%d): %d", constants.MaxPort, start)
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid end port in range %s: %v", pr, err)
	}
	if end < 0 || end > constants.MaxPort {
		return nil, fmt.Errorf("end port out of range (0-%d): %d", constants.MaxPort, end)
	}
	if start > end {
		return nil, fmt.Errorf("start port cannot be greater than end port in range %s", pr)
	}

	var ports []int
	for port := start; port <= end; port++ {
		ports = append(ports, port)
	}

	return ports, nil
}

// TODO doc
func ParseSinglePort(pr string) (int, error) {
	port, err := strconv.Atoi(pr)
	if err != nil {
		return 0, fmt.Errorf("invalid port number: %s", pr)
	}
	if port < 0 || port > constants.MaxPort {
		return 0, fmt.Errorf("port out of range (0-%d): %d", constants.MaxPort, port)
	}

	return port, nil
}

// TODO doc
func DeduplicatePorts(ports []int) []int {
	keys := make(map[int]bool)
	var uniquePorts []int

	for _, port := range ports {
		if _, value := keys[port]; !value {
			keys[port] = true
			uniquePorts = append(uniquePorts, port)
		}
	}

	return uniquePorts
}
