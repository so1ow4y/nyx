package config

import (
	"testing"
	"time"
)

func TestDefaults(t *testing.T) {
	if DefaultPorts != "80,443" {
		t.Errorf("Expected DefaultPorts to be '80,443', got '%s'", DefaultPorts)
	}

	if DefaultProtocol != "tcp" {
		t.Errorf("Expected DefaultProtocol to be 'tcp', got '%s'", DefaultProtocol)
	}

	if DefaultTimeout != 2*time.Second {
		t.Errorf("Expected DefaultTimeout to be 2s, got %v", DefaultTimeout)
	}
}
