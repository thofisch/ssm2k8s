package k8s

import (
	"testing"
)

func Test_verifyVersion(t *testing.T) {
	serverVersion := "v0.3.5"

	tests := []string{
		"v0.3.0",
		"0.3.0",
		"0.3.99",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			if err := verifyVersion(serverVersion, tt); err != nil {
				t.Errorf("verifyVersion() error = %v", err)
			}
		})
	}
}
