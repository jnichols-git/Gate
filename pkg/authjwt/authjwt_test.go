package authjwt

import (
	"strings"
	"testing"
)

func TestNewJWT(t *testing.T) {
	testUser := "testUser@gmail.com"
	testPerm := map[string]bool{
		"authorized": true,
	}
	token := NewJWT(testUser, testPerm)
	// Verify that all token fields are properly set
	if token.Header.Algorithm != "sha256" {
		t.Errorf("Token uses unrecognized algorithm %s", token.Header.Algorithm)
	}
	if token.Header.Type != "jwt" {
		t.Errorf("Token header has incorrect type %s", token.Header.Type)
	}
	if token.Body.ForUser != testUser {
		t.Errorf("Token has incorrect user %s", token.Body.ForUser)
	}
	if !token.Body.Access["authorized"] {
		t.Errorf("Token perm authorization has incorrect value %t", token.Body.Access["authorization"])
	}
}

// Test Export/Verify success cases. This should run without errors.
func TestExportVerifyValid(t *testing.T) {
	testUser := "testUser@gmail.com"
	testPerm := map[string]bool{
		"authorized": true,
	}
	jwt := NewJWT(testUser, testPerm)
	token := Export(jwt, []byte("test"))
	res, valid, err := Verify(token, []byte("test"))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if !valid {
		t.Logf("Token didn't validate.")
		if jwt.Body.ForUser != res.Body.ForUser {
			t.Errorf("Original token ForUser %s != result token ForUser %s.", jwt.Body.ForUser, res.Body.ForUser)
		}
		if jwt.Body.Access["user-type"] != res.Body.Access["user-type"] {
			t.Errorf("Original token authorization %t != result token authorization %t.", jwt.Body.Access["authorization"], res.Body.Access["authorization"])
		}
		t.Errorf("Fields validated, validation error unknown.")
	}
}

// Test Export/Verify failure cases. This should catch an error on *almost* every call.
func TestExportVerifyInvalid(t *testing.T) {
	testUser := "testUser@gmail.com"
	testPerm := map[string]bool{
		"authorized": true,
	}
	jwt := NewJWT(testUser, testPerm)
	Export(jwt, []byte("test"))
	// Create an actual valid token
	jwt = NewJWT(testUser, testPerm)
	token := Export(jwt, []byte("test"))
	// Testing basic modifications to token in header and body

	// Header
	headerIndex := 0
	invalidHeaderToken := []rune(token)
	// Invalid base64 rune
	invalidHeaderToken[headerIndex] = '='
	_, _, err := Verify(string(invalidHeaderToken), []byte("test"))
	if err == nil {
		t.Error("Placing invalid rune in header didn't return an error from Verify")
	}
	invalidHeaderToken[headerIndex] = 'a'
	_, _, err = Verify(string(invalidHeaderToken), []byte("test"))
	if err == nil {
		t.Error("Placing erroneous rune in header didn't return an error from Verify")
	}

	// Body
	bodyIndex := strings.Index(token, ".") + 1
	invalidBodyToken := []rune(token)
	// Invalid base64 rune
	invalidBodyToken[bodyIndex] = '='
	_, _, err = Verify(string(invalidBodyToken), []byte("test"))
	if err == nil {
		t.Error("Placing invalid rune in body didn't return an error from Verify")
	}
	invalidBodyToken[bodyIndex] = 'a'
	_, _, err = Verify(string(invalidBodyToken), []byte("test"))
	if err == nil {
		t.Error("Placing erroneous rune in body didn't return an error from Verify")
	}
}
