package param

import (
	"fmt"
	"reflect"
	"testing"
)

func assertEmpty(t *testing.T, i interface{}) {
	ok, length := getLength(i)

	if !ok {
		_type := reflect.TypeOf(i)

		t.Fatalf("\"%s\" could not be applied builtin len()", _type.Name())
		return
	} else if length != 0 {
		t.Fatalf("Expected empty but was %v", i)
	}
}

func getLength(i interface{}) (ok bool, length int) {
	defer func() {
		if e := recover(); e != nil {
			ok = false
		}
	}()

	v := reflect.ValueOf(i)

	return true, v.Len()
}

func assertTrue(t *testing.T, b bool) {
	if !b {
		t.Errorf("Expected 'true'")
	}
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected == actual {
		return
	}
	message := fmt.Sprintf("%v != %v", expected, actual)
	t.Fatal(message)
}

func assertOk(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected OK, but got: %s", err)
	}
}

func assertNotOk(t *testing.T, err error) {
	if err == nil {
		t.Fatalf("Expected not OK, but got: OK")
	}
}
