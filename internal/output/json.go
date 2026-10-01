package output

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlab.com/ere6u5/nyx/pkg/types"
)

// TODO doc
func SaveToJSONFile(filename string, data []types.ScanResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
