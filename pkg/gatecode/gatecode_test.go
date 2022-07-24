package gatecode

import "testing"

func TestValidateAuthCode(t *testing.T) {
	// Test correct authcode
	code := NewGateCode("testuser@gmail.com")
	valid := ValidateGateCode("testuser@gmail.com", code)
	if !valid {
		t.Errorf("Validation 0 failed when it should have passed")
	}
	// Test incorrect gate code
	code = NewGateCode("testuser@gmail.com")
	valid = ValidateGateCode("testuser@gmail.com", "AAAAA")
	if valid {
		t.Errorf("Validation 1 passed when it should have failed")
	}
	// Test correct gate code with incorrect email
	code = NewGateCode("testuser@gmail.com")
	valid = ValidateGateCode("otheruser@gmail.com", code)
	if valid {
		t.Errorf("Validation 2 passed when it should have failed")
	}
}
