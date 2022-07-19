package database

import (
	"encoding/json"

	"github.com/pkg/errors"
)

const testing bool = true

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
	authDBLayout, err = authDBLayout.opened(!testing)
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
	salt, err := genSalt()
	if err != nil {
		return err
	}
	// Check to make sure we're not re-registering a user
	if registered, err := keyExists(authDBLayout["users"], authDBKey(email, username)); err != nil {
		return err
	} else {
		if registered {
			return errors.New("User already registered")
		}
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
	return username == entry.Credentials.Username && pwdHash == entry.Credentials.PasswordHash, nil
}
