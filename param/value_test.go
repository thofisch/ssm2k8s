package param

import (
	"fmt"
	"testing"
)

func TestParameterValueSecretString(t *testing.T) {
	assertEqual(t, "*****", ParameterValueSecretString)
}

func TestParameterValue_String(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		secret   bool
	}{
		{name: "secret", expected: ParameterValueSecretString, secret: true},
		{name: "clear", expected: "val", secret: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pv := NewParameterValue("val", test.secret)

			assertEqual(t, test.expected, fmt.Sprintf("%s", pv))
		})
	}
}

func TestParameterValue_GetValue(t *testing.T) {
	pv := NewParameterValue("val", true)

	assertEqual(t, "val", pv.GetValue())
}

func TestParameterValue_IsSecret(t *testing.T) {
	tests := []struct {
		name     string
		secret   bool
	}{
		{name: "secret", secret: true},
		{name: "clear", secret: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pv := NewParameterValue("val", test.secret)

			assertEqual(t, test.secret, pv.IsSecret())
		})
	}}

