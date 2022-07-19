package database

import (
	"testing"
)

func TestDBAccess(t *testing.T) {
	DATABASE_TESTING = true
	// Open database
	err := OpenDB()
	if err != nil {
		t.Error(err)
	}
	// Close database when done.
	defer CloseDB()
	if valid, err := ValidateUserCred("user@email.com", "username", "password"); valid || err == nil {
		t.Error("Validated user credentials that haven't been registered")
	}

	// Register a user
	if err := RegisterUser("user@email.com", "username", "password", nil); err != nil {
		t.Error(err)
	}
	// Validate those user credentials
	if _, err := ValidateUserCred("user@email.com", "username", "password"); err != nil {
		t.Error(err)
	}
	// Validate again, but using the wrong password
	if valid, _ := ValidateUserCred("user@email.com", "username", "wrongpassword"); valid {
		t.Error("Validated user credentials that are incorrect")
	}
	// Try re-registering (should fail)
	if err := RegisterUser("user@email.com", "username", "password", nil); err == nil {
		t.Error("Re-registration under same email/username succeeded when it should fail")
	}
	// Given that fails, let's change passwords. Wrong username?
	if err := ChangeUserPassword("user@email.com", "wrongusername", "password", "newPassword"); err == nil {
		t.Error("Password change with incorrect username succeeded when it should fail")
	}
	// Wrong old password
	if err := ChangeUserPassword("user@email.com", "username", "wrongpassword", "newPassword"); err == nil {
		t.Error("Password change with incorrect old password succeeded when it should fail")
	}
	// Correct everything
	if err := ChangeUserPassword("user@email.com", "username", "password", "newPassword"); err != nil {
		t.Error(err)
	}
	// Re-validate with new password
	if _, err := ValidateUserCred("user@email.com", "username", "newPassword"); err != nil {
		t.Error(err)
	}
}
