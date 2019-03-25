package assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)


func Equal(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texpected: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, expected, actual)
		tb.FailNow()
	}
}

func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func NotOk(tb testing.TB, err error) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: expected error didn't occur\033[39m\n\n", filepath.Base(file), line)
		tb.FailNow()
	}
}

func True(tb testing.TB, b bool) {
	if !b {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: expected truen\n\tgot: false\033[39m\n\n", filepath.Base(file), line)
		tb.FailNow()
	}
}

//func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
//	if !condition {
//		_, file, line, _ := runtime.Caller(1)
//
//		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
//		tb.FailNow()
//	}
//}

