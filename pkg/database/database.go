package database

import (
	"encoding/json"
	"strings"

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

func authDBUnKey(key []byte) (string, string) {
	emun := strings.Split(string(key), ".")
	un, _ := stringDecode(emun[0])
	em, _ := stringDecode(emun[1])
	return string(un), string(em)
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
	if EmailTaken(email) {
		return errors.New("Email is already in use.")
	}
	if UsernameTaken(username) {
		return errors.New("Username is already in use.")
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

// Check if an email is taken
func EmailTaken(email string) bool {
	key, _ := findUserByEmail(authDBLayout["users"], email)
	return key != nil
}

// Check if a username is taken
func UsernameTaken(username string) bool {
	key, _ := findUserByUsername(authDBLayout["users"], username)
	return key != nil
}

func UserExists(email, username string) (bool, error) {
	return keyExists(authDBLayout["users"], authDBKey(email, username))
}

func ValidateUserCred(username, password string) (bool, UserEntry, error) {
	key, err := findUserByUsername(authDBLayout["users"], username)
	if err != nil || key == nil {
		return false, UserEntry{}, err
	}
	dbval, err := getKey(authDBLayout["users"], key)
	if err != nil {
		return false, UserEntry{}, err
	}
	entry := UserEntry{}
	if err := json.Unmarshal(dbval, &entry); err != nil {
		return false, UserEntry{}, err
	}
	salt, _ := stringDecode(entry.Credentials.PasswordSalt)
	pwdHash, err := slowHash([]byte(password), salt, entry.Credentials.HashFunc)
	if err != nil {
		return false, UserEntry{}, err
	}
	valid := username == entry.Credentials.Username && pwdHash == entry.Credentials.PasswordHash
	if valid {
		return valid, entry, nil
	} else {
		return valid, UserEntry{}, nil
	}
}

func ChangeUserPassword(username, password string, newPassword string) error {
	oldOK, entry, _ := ValidateUserCred(username, password)
	// Need to validate old password to change to a new one.
	if !oldOK {
		return errors.New("Password change failed: old password incorrect")
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
	err = setKey(authDBLayout["users"], authDBKey(entry.Credentials.Email, username), out)
	if err != nil {
		return err
	}
	return nil
}
