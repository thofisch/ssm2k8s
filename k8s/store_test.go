package k8s

import (
	"testing"
)

func Test_verifyVersion(t *testing.T) {
	tests := []string{
		"v1.2.3",
		"1.2.3",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			if err := verifyVersion(tt); err != nil {
				t.Errorf("verifyVersion() error = %v", err)
			}
		})
	}
}
