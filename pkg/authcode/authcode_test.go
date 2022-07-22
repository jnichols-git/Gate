package authcode

import "testing"

func TestValidateAuthCode(t *testing.T) {
	// Test correct authcode
	code := NewAuthCode("testuser@gmail.com")
	valid := ValidateAuthCode("testuser@gmail.com", code)
	if !valid {
		t.Errorf("Validation 0 failed when it should have passed")
	}
	// Test incorrect authcode
	code = NewAuthCode("testuser@gmail.com")
	valid = ValidateAuthCode("testuser@gmail.com", "AAAAA")
	if valid {
		t.Errorf("Validation 1 passed when it should have failed")
	}
	// Test correct authcode with incorrect email
	code = NewAuthCode("testuser@gmail.com")
	valid = ValidateAuthCode("otheruser@gmail.com", code)
	if valid {
		t.Errorf("Validation 2 passed when it should have failed")
	}
}
