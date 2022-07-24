package credentials

import (
	"os"
	"testing"
)

func TestDBAccess(t *testing.T) {
	// Delete previous test database
	err := os.Remove("test.db")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// Every single access function should fail if the database is not open.
	if eAdd, eUpdate := addUser(nil), updateUser(nil); eAdd == nil || eUpdate == nil {
		t.Error("addUser and updateUser should both return error with a closed database")
	}
	if _, eFindEmail := findUserEntryByEmail("foo@bar.com"); eFindEmail == nil {
		t.Error("findUserEntryByEmail should return error with a closed database")
	}
	if _, eFindUname := findUserEntryByUsername("foobar"); eFindUname == nil {
		t.Error("findUserEntryByUsername should return error with a closed database")
	}

	// Open database
	if err := OpenDB("test.db"); err != nil {
		t.Error(err)
	}

	if valid, invalidUser, _ := ValidateUserCred("username", "password"); valid {
		t.Error("Validated user credentials that haven't been registered")
	} else {
		if !invalidUser.Empty() {
			t.Error("ValidateUserCred failed, but the returned User was not the empty User")
		}
	}
	if invalidUser, err := findUserEntryByUsername("username"); err != nil {
		t.Error("Standalone findUserEntryByUsername returned no error for a nonexistent user")
	} else {
		if !invalidUser.Empty() {
			t.Error("findUserEntryByUsername failed, but the returned userEntry was not the empty userEntry")
		} else if !invalidUser.toUser().Empty() {
			t.Error("Empty userEntry was returned, but did not convert to the empty User with toUser()")
		}
	}

	// Register a user
	if err := RegisterUser("user@email.com", "username", "password", nil); err != nil {
		t.Error(err)
	}

	// With that user registered, we should be able to find them by email or username...
	if foundUser, err := FindUserByEmail("user@email.com"); err != nil {
		t.Error("Got an error when finding user@email.com after registration.")
	} else if foundUser.Empty() {
		t.Error("FindUserByEmail returned the empty user for user@email.com after registration.")
	}
	if foundUser, err := FindUserByUsername("username"); err != nil {
		t.Error("Got an error when finding username after registration.")
	} else if foundUser.Empty() {
		t.Error("FindUserByUsername returned the empty user username after registration.")
	}
	// There should also be 1 entry.
	// This doesn't actually hold for the prod db, since there's an admin account, but we're going to ignore that for this test.
	if Entries() != 1 {
		t.Errorf("Databse should have 1 entry, but has %d instead", Entries())
	}

	// Validate those user credentials
	if _, _, err := ValidateUserCred("username", "password"); err != nil {
		t.Error(err)
	}

	// Validate again, but using the wrong password
	if valid, _, _ := ValidateUserCred("username", "wrongpassword"); valid {
		t.Error("Validated user credentials that are incorrect")
	}

	// Try re-registering (should fail)
	if err := RegisterUser("user@email.com", "username", "password", nil); err == nil {
		t.Error("Re-registration under same email/username succeeded when it should fail")
	}
	// Make sure the number of entries didn't change with a failed re-registration
	if Entries() != 1 {
		t.Errorf("Database should have 1 entry, but has %d instead", Entries())
	}

	// Given that fails, let's change passwords. Wrong username?
	if err := ChangeUserPassword("wrongusername", "password", "newPassword"); err == nil {
		t.Error("Password change with incorrect username succeeded when it should fail")
	}

	// Wrong old password
	if err := ChangeUserPassword("username", "wrongpassword", "newPassword"); err == nil {
		t.Error("Password change with incorrect old password succeeded when it should fail")
	}

	// Correct everything
	if err := ChangeUserPassword("username", "password", "newPassword"); err != nil {
		t.Error(err)
	}

	// Re-validate with new password
	if _, _, err := ValidateUserCred("username", "newPassword"); err != nil {
		t.Error(err)
	}

	// Change the permissions for this user
	if err := ChangeUserPermissions("username", map[string]bool{"testPermission": true}); err != nil {
		t.Error(err)
	}
	// Check that the permissions changed.
	if user, err := FindUserByUsername("username"); err != nil || !user.Permissions["testPermission"] {
		t.Error("Failed to set testPermission; getting user afterwards did not reflect the permission change.")
	}
}
