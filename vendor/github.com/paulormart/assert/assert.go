// Assert
// ideas borrowed from https://github.com/stretchr/testify
package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// getTestFile identifies the test file that has been called highlighting the line
// Note: we assume a go test file always ends as _test.go
func getTestFile() string {
	for i := 0; ; i++ {
		_, file, line, ok := runtime.Caller(i)

		if !ok || file == "" || line == 0 {
			break
		}

		if x := strings.LastIndex(file, "/"); x != -1 {
			if name := file[x:]; strings.Contains(name, "_test.go") {
				return fmt.Sprintf(".%s:%d", name, line)
			}
		}
	}
	return ""
}

// Fail get Caller file and print assertion
func Fail(t *testing.T, expected, actual interface{}) bool {
	if caller := getTestFile(); caller != "" {
		t.Errorf("\rError at %v\t Not equal: %v (expected) != %v", caller, expected, actual)
	}
	return false
}

// Equal asserts deep equal for two objects
func Equal(t *testing.T, expected, actual interface{}) bool {

	if expected == nil || actual == nil {
		return expected == actual
	}

	if !reflect.DeepEqual(expected, actual) {
		return Fail(t, expected, actual)
	}

	return true
}
