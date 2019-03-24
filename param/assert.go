package param

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		return
	}

	t.Fatalf("Expected: %v, got: %v", expected, actual)
}

func assertOk(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected OK, but got: %s", err)
	}
}

func assertNotOk(t *testing.T, err error) {
	if err == nil {
		t.Fatalf("Expected error, but got: OK")
	}
}

func assertTrue(t *testing.T, b bool) {
	if !b {
		t.Errorf("Expected 'true'")
	}
}
