package authcode

import "testing"

func TestNewAuthCode(t *testing.T) {
	code := NewAuthCode("test-user")
	if code.ForUser != "test-user" {
		t.Errorf("Expected user %s, got %s", "test-user", code.ForUser)
	}
}

func TestValidateAuthCode(t *testing.T) {
	code := NewAuthCode("test-user")
	valid := ValidateAuthCode("test-user", code.Code)
	if !valid {
		t.Errorf("Validation 0 failed when it should have passed")
	}
	code = NewAuthCode("test-user")
	valid = ValidateAuthCode("test-user", "AAAAA")
	if valid {
		t.Errorf("Validation 1 passed when it should have failed")
	}
	code = NewAuthCode("test-user")
	valid = ValidateAuthCode("other-user", code.Code)
	if valid {
		t.Errorf("Validation 2 passed when it should have failed")
	}
}
