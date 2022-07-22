package authcred

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/* From config.go: DBConfig
Path string
*/

type UserCred struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserPerm map[string]bool

type User struct {
	Credentials UserCred `json:"credentials"`
	Permissions UserPerm `json:"permissions"`
}

// userEntry: User representation of the database
type userEntry struct {
	ID           uint   `gorm:"autoIncrement,primaryKey"`
	Email        string `gorm:"email"`
	Username     string `gorm:"username"`
	PasswordHash string `gorm:"password"`
	Salt         string `gorm:"salt"`
	HashFunc     string `gorm:"hashfunc"`
	Permissions  string `gorm:"permissions"`
}

func (u userEntry) ToUser() User {
	outUser := User{
		Credentials: UserCred{
			u.Email,
			u.Username,
		},
		Permissions: make(UserPerm),
	}
	json.Unmarshal([]byte(u.Permissions), &outUser.Permissions)
	return outUser
}

// Open the database
func OpenDB(path string) error {
	var err error
	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&userEntry{})
	if err != nil {
		return err
	}

	return nil
}

func Entries() int {
	entries := make([]userEntry, 0)
	db.Find(&entries)
	return len(entries)
}

// Exported version of findUserByEmail; returns public User instead of userEntry
func FindUserByEmail(email string) (User, error) {
	uentry, err := findUserByEmail(email)
	if err != nil {
		return User{}, err
	}
	return uentry.ToUser(), nil
}

// Register a user with credentials and permissions
// Returns error or nil
func RegisterUser(email, username, password string, permissions UserPerm) error {
	if user, _ := findUserByEmail(email); user.Email != "" {
		return errors.New("Email is already in use.")
	}
	if user, _ := findUserByUsername(username); user.Username != "" {
		return errors.New("Username is already in use.")
	}
	perm, err := json.Marshal(permissions)
	salt, err := genSalt()
	if err != nil {
		return err
	}
	pwdHash, err := slowHash([]byte(password), salt, "sha512")
	if err != nil {
		return err
	}
	entry := &userEntry{
		Email:        email,
		Username:     username,
		PasswordHash: pwdHash,
		Salt:         stringEncode(salt),
		HashFunc:     "sha512",
		Permissions:  string(perm),
	}
	addUser(entry)
	return nil
}

// Validate a user with credentials
// Returns success, user info, and error/nil
func ValidateUserCred(username, password string) (bool, *User, error) {
	// Find user. Fail out if non-eistent
	user, err := findUserByUsername(username)
	if err != nil {
		return false, nil, err
	}
	if user.Username == "" {
		return false, nil, fmt.Errorf("User %s not found", username)
	}
	// Check hashed password against input password.
	salt, _ := stringDecode(user.Salt)
	pwdHash, err := slowHash([]byte(password), salt, user.HashFunc)
	if err != nil {
		return false, nil, err
	}
	valid := username == user.Username && pwdHash == user.PasswordHash
	if valid {
		outUser := &User{
			Credentials: UserCred{
				user.Email,
				user.Username,
			},
			Permissions: make(UserPerm),
		}
		json.Unmarshal([]byte(user.Permissions), &outUser.Permissions)
		return true, outUser, nil
	} else {
		return false, &User{}, fmt.Errorf("User validation failed")
	}
}

// Change user password
// Returns success and error/nil
func ChangeUserPassword(username, password string, newPassword string) error {
	valid, _, err := ValidateUserCred(username, password)
	if !valid {
		if err != nil {
			return err
		} else {
			return fmt.Errorf("User validation failed")
		}
	}
	user, err := findUserByUsername(username)
	if err != nil {
		return err
	}
	// If the user didn't exist, validation would have failed
	salt, _ := genSalt()
	pwdHash, err := slowHash([]byte(newPassword), salt, user.HashFunc)
	if err != nil {
		return err
	}
	user.PasswordHash = pwdHash
	user.Salt = stringEncode(salt)
	updateUser(&user)
	return nil
}
