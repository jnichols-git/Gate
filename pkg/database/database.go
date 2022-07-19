package database

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

var DATABASE_TESTING bool = true

type UserCred struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
	PasswordSalt string `json:"salt"`
	HashFunc     string `json:"hf"`
}

type UserPerm map[string]bool

type UserEntry struct {
	Credentials UserCred `json:"credentials"`
	Permissions UserPerm `json:"permissions"`
}

var authDBLayout DatabaseCollection = DatabaseCollection{
	"users": nil,
}

func authDBKey(email, username string) []byte {
	email = stringEncode([]byte(email))
	username = stringEncode([]byte(username))
	return []byte(email + "." + username)
}

func OpenDB() error {
	var err error
	authDBLayout, err = authDBLayout.opened(!DATABASE_TESTING)
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() error {
	for _, db := range authDBLayout {
		err := db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterUser(email, username, password string, permissions UserPerm) error {
	// Check to make sure we're not re-registering a user
	if registered, err := keyExists(authDBLayout["users"], authDBKey(email, username)); err != nil {
		return err
	} else {
		if registered {
			return errors.New("User already registered")
		}
	}
	salt, err := genSalt()
	if err != nil {
		return err
	}
	pwdHash, err := slowHash([]byte(password), salt, "sha512")
	if err != nil {
		return err
	}
	entry := UserEntry{
		Credentials: UserCred{
			email,
			username,
			pwdHash,
			stringEncode(salt),
			"sha512",
		},
		Permissions: permissions,
	}
	out, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	err = setKey(authDBLayout["users"], authDBKey(email, username), out)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUserCred(email, username, password string) (bool, error) {
	dbval, err := getKey(authDBLayout["users"], authDBKey(email, username))
	if err != nil {
		return false, err
	}
	entry := UserEntry{}
	if err := json.Unmarshal(dbval, &entry); err != nil {
		return false, err
	}
	salt, _ := stringDecode(entry.Credentials.PasswordSalt)
	pwdHash, err := slowHash([]byte(password), salt, entry.Credentials.HashFunc)
	if err != nil {
		return false, err
	}
	fmt.Printf("%s:\n%s\n%s\n", username, pwdHash, entry.Credentials.PasswordHash)
	fmt.Printf("%t, %t\n", username == entry.Credentials.Username, pwdHash == entry.Credentials.PasswordHash)
	return username == entry.Credentials.Username && pwdHash == entry.Credentials.PasswordHash, nil
}

func ChangeUserPassword(email, username, password string, newPassword string) error {
	// Get current user data. This will fail if the user does not exist.
	dbval, err := getKey(authDBLayout["users"], authDBKey(email, username))
	if err != nil {
		return err
	}
	entry := UserEntry{}
	if err := json.Unmarshal(dbval, &entry); err != nil {
		return err
	}
	// Need to validate old password to change to a new one.
	if oldOk, _ := ValidateUserCred(email, username, password); !oldOk {
		return errors.New("Password change failed: old password incorrect")
	} else {
	}
	// Generate password hash
	salt, err := genSalt()
	if err != nil {
		return err
	}
	pwdHash, err := slowHash([]byte(newPassword), salt, "sha512")
	if err != nil {
		return err
	}
	entry.Credentials.PasswordHash = pwdHash
	entry.Credentials.PasswordSalt = stringEncode(salt)
	// Write out new entry
	out, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	err = setKey(authDBLayout["users"], authDBKey(email, username), out)
	if err != nil {
		return err
	}
	return nil
}
