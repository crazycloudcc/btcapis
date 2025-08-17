// Package assert 提供断言工具
package assert

import (
	"fmt"
	"testing"
)

// Equal 断言两个值相等
func Equal(t *testing.T, expected, actual interface{}, msg ...string) {
	t.Helper()
	if expected != actual {
		message := "values not equal"
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf("%s: expected %v, got %v", message, expected, actual)
	}
}

// NotEqual 断言两个值不相等
func NotEqual(t *testing.T, expected, actual interface{}, msg ...string) {
	t.Helper()
	if expected == actual {
		message := "values should not be equal"
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf("%s: expected not %v, got %v", message, expected, actual)
	}
}

// Nil 断言值为nil
func Nil(t *testing.T, actual interface{}, msg ...string) {
	t.Helper()
	if actual != nil {
		message := "value should be nil"
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf("%s: expected nil, got %v", message, actual)
	}
}

// NotNil 断言值不为nil
func NotNil(t *testing.T, actual interface{}, msg ...string) {
	t.Helper()
	if actual == nil {
		message := "value should not be nil"
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf("%s: expected not nil, got nil", message)
	}
}

// True 断言值为true
func True(t *testing.T, actual bool, msg ...string) {
	t.Helper()
	if !actual {
		message := "value should be true"
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf("%s: expected true, got false", message)
	}
}

// False 断言值为false
func False(t *testing.T, actual bool, msg ...string) {
	t.Helper()
	if actual {
		message := "value should be false"
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf("%s: expected false, got true", message)
	}
}

// Error 断言有错误
func Error(t *testing.T, actual error, msg ...string) {
	t.Helper()
	if actual == nil {
		message := "expected error, got nil"
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf(message)
	}
}

// NoError 断言没有错误
func NoError(t *testing.T, actual error, msg ...string) {
	t.Helper()
	if actual != nil {
		message := "expected no error"
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf("%s: got error %v", message, actual)
	}
}

// Contains 断言字符串包含子串
func Contains(t *testing.T, str, substr string, msg ...string) {
	t.Helper()
	if !contains(str, substr) {
		message := fmt.Sprintf("string should contain %q", substr)
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf("%s: %q does not contain %q", message, str, substr)
	}
}

// NotContains 断言字符串不包含子串
func NotContains(t *testing.T, str, substr string, msg ...string) {
	t.Helper()
	if contains(str, substr) {
		message := fmt.Sprintf("string should not contain %q", substr)
		if len(msg) > 0 {
			message = msg[0]
		}
		t.Errorf("%s: %q contains %q", message, str, substr)
	}
}

// contains 检查字符串是否包含子串
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(substr) == 0 ||
		(len(str) > len(substr) && (str[:len(substr)] == substr ||
			str[len(str)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(str)-len(substr); i++ {
					if str[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
