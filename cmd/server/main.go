package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"gitlab.com/ere6u5/nyx/internal/core"
	"gitlab.com/ere6u5/nyx/internal/parser"
	"gitlab.com/ere6u5/nyx/pkg/types"
)

func main() {
	http.Handle("/", http.FileServer((http.Dir("cmd/server/static"))))

	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req types.ScanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		ports, err := parser.ParsePorts(req.Ports)
		if err != nil {
			http.Error(w, "Invalid ports format", http.StatusBadRequest)
		}

		timeout, err := ParseDuration(req.Timeout)
		if err != nil {
			http.Error(w, "Invalid timeout format: "+err.Error(), http.StatusBadRequest)
		}

		result := core.Scan(req.Host, ports, req.Protocol, timeout)

		renderResults(w, result)
	})

	log.Println("Server stated at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO doc
func ParseDuration(input string) (time.Duration, error) {
	if input == "" {
		return 2 * time.Second, nil
	}

	if val, err := strconv.Atoi(input); err == nil {
		return time.Duration(val) * time.Second, nil
	}

	re := regexp.MustCompile(`^(\d+)(ms|s|m|h)$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return 0, errors.New("invalid duration format. Use: '300ms', '2s', '1m', '1h'")
	}

	val, _ := strconv.Atoi(matches[1])
	unit := matches[2]

	switch unit {
	case "ms":
		return time.Duration(val) * time.Millisecond, nil
	case "s":
		return time.Duration(val) * time.Second, nil
	case "m":
		return time.Duration(val) * time.Minute, nil
	case "h":
		return time.Duration(val) * time.Hour, nil
	default:
		return 0, errors.New("unknown time unit")
	}
}

// TODO doc
func renderResults(w http.ResponseWriter, results []types.ScanResult) {
	tmpl := template.Must(template.ParseFiles("cmd/server/static/result.html"))
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, results)
}
