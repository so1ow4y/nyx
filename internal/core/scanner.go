package core

import (
	"fmt"
	"net"
	"sync"
	"time"

	"gitlab.com/ere6u5/nyx/pkg/types"
)

// TODO doc
func Scan(host string, ports []int, protocol string, timeout time.Duration) []types.ScanResult {

	var (
		results []types.ScanResult
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	for _, pr := range ports {
		wg.Add(1)
		go func(goPort int) {
			defer wg.Done()
			conn, err := net.DialTimeout(protocol, fmt.Sprintf(host+":%d", goPort), timeout)

			result := types.ScanResult{Port: goPort}

			if err != nil {
				result.Status = "closed"
			} else {
				result.Status = "open"
				conn.Close()
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(pr)
	}

	wg.Wait()
	return results
}
