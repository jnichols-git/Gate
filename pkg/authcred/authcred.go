// Package authcred handles the authentication of users using username-password pairs
// Each user is stored as a userEntry in a local database. The database columns are defined
// by the userEntry struct, so they appear as seen below:
//   +----+-------+----------+---------------+------+---------------+-------------+
//   | ID | Email | Username | Password Hash | Salt | Hash Function | Permissions |
//   +----+-------+----------+---------------+------+---------------+-------------+
// Two main authentication functions are provided in RegisterUser() and ValidateUserCred(),
// with supporting functions ChangeUserPassword() and ChangeUserPermissions() to alter the
// data of already-existing users in the database.
//
// authcred exports the User type, which contains the same data as userEntry with private data
// (password hash, salt, hash func, internal ID) removed. ValidateUserCred() returns one, so that
// when the authentication API is called, it returns back information about the user in a format
// that can easily pass back to the application servers or converted into a token without exposing
// important data.
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

// A User contains *public* information about a user.
// authcred functions that return user info MUST return this.
type User struct {
	Email       string          `json:"email"`
	Username    string          `json:"username"`
	Permissions map[string]bool `json:"permissions"`
}

// Empty() checks if the calling User is the empty User.
// This means no email or username, and 0 permissions.
//
// Calling:
//   - u User: User to check.
// Output:
//   - bool: Is u the empty User?
func (u User) Empty() bool {
	return u.Email == "" && u.Username == "" && len(u.Permissions) == 0
}

// A userEntry contains public *and private* information about a user.
// authcred functions that are exported MUST NOT return this.
type userEntry struct {
	ID           uint   `gorm:"autoIncrement,primaryKey"`
	Email        string `gorm:"email"`
	Username     string `gorm:"username"`
	PasswordHash string `gorm:"password"`
	Salt         string `gorm:"salt"`
	HashFunc     string `gorm:"hashfunc"`
	Permissions  string `gorm:"permissions"`
}

// Empty() checks if the calling userEntry is the empty userEntry.
//
// Calling:
//   - u userEntry: userEntry to check.
// Output:
//   - bool: Is u the empty userEntry?
func (u userEntry) Empty() bool {
	return u.ID == 0 && u.Email == "" && u.Username == "" && u.PasswordHash == "" && u.Salt == "" && u.HashFunc == "" && u.Permissions == ""
}

// toUser converts the calling userEntry into a User.
//
// Calling:
//   - u userEntry: Data to convert. Notably, the zero userEntry converts to the zero User.
// Output:
//   - User: the resulting public User struct.
func (u userEntry) toUser() User {
	outUser := User{
		Email:       u.Email,
		Username:    u.Username,
		Permissions: make(map[string]bool),
	}
	json.Unmarshal([]byte(u.Permissions), &outUser.Permissions)
	return outUser
}

// Open the database.
// This MUST be called before any authcred operations take place, and if the path is changed,
// a different DB will be opened; this is configured under DB.Path in config.yml, and should probably not
// change unless you have a testing database to use.
// Only one database can be open at a time.
//
// Input:
//   - path string: Path to the database file
// Output:
//   - error: Output if the open fails, or if the userEntry struct changes in a way that prevents migrating the DB.
//   This most commonly occurs if the path does not exist; gorm can create a new file, but not directories.
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

// Get the current number of entries in the database.
// It's not likely that this has significant use outside of noticing if the user is initializing a new database; to preserve
// security in this case, the user should be asked to create the first account as an admin account before opening to a network.
//
// Output:
//   - int: Number of entries in the current database
func Entries() int {
	entries := make([]userEntry, 0)
	db.Find(&entries)
	return len(entries)
}

// Exported version of findUserEntryByEmail; returns public User instead of userEntry.
//
// Input:
//   - email string: Email to find
// Output:
//   - User: User data, or empty user if not found.
func FindUserByEmail(email string) (User, error) {
	uentry, err := findUserEntryByEmail(email)
	if err != nil {
		return User{}, err
	}
	return uentry.toUser(), nil
}

// Exported version of findUserEntryByUsername; returns public User instead of userEntry.
//
// Input:
//   - username string: Username to find
// Output:
//   - User: User data, or empty user if not found.
func FindUserByUsername(username string) (User, error) {
	uentry, err := findUserEntryByUsername(username)
	if err != nil {
		return User{}, err
	}
	return uentry.toUser(), nil
}

// Register a user with the given credentials and permissions.
//
// Input:
//   - email, username string: User email/username pair. Both of these values MUST be unique.
//   - password string: User password. To avoid conflicts with integration, auth imposes no password restrictions; it is
//   expected that the application manage restrictions such as password length.
//   - permissions map[string]bool: User permissions. auth only takes advantage of the admin permission; all others are
//   application-defined.
// Output:
//   - error: Any errors that occur during the registration of a user, including: non-unique email/username, failure to generate
//   password salt, failure to hash password. If an error is returned, no change is made to the database.
func RegisterUser(email, username string, password string, permissions map[string]bool) error {
	if user, _ := findUserEntryByEmail(email); user.Email != "" {
		return errors.New("Email is already in use.")
	}
	if user, _ := findUserEntryByUsername(username); user.Username != "" {
		return errors.New("Username is already in use.")
	}
	perm, _ := json.Marshal(permissions)
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

// Validate a user with username and password credentials.
//
// Input:
//   - username, password string: User credentials. The username will be used to find the userEntry, and then the password
//   will be hashed with that userEntry salt and hash function and compared to the userEntry password.
// Output:
//   - bool: Is user valid?
//   - User: Public user credentials.
//   - error: Any errors that occur during user validation, including: failure to find user, failure to hash password, failure to validate
//   user.
func ValidateUserCred(username, password string) (bool, User, error) {
	// Find user. Fail out if non-eistent
	user, err := findUserEntryByUsername(username)
	if err != nil {
		return false, User{}, err
	}
	if user.Username == "" {
		return false, User{}, fmt.Errorf("User %s not found", username)
	}
	// Check hashed password against input password.
	salt, _ := stringDecode(user.Salt)
	pwdHash, err := slowHash([]byte(password), salt, user.HashFunc)
	if err != nil {
		return false, User{}, err
	}
	valid := username == user.Username && pwdHash == user.PasswordHash
	if valid {
		outUser := User{
			Email:       user.Email,
			Username:    user.Username,
			Permissions: make(map[string]bool),
		}
		json.Unmarshal([]byte(user.Permissions), &outUser.Permissions)
		return true, outUser, nil
	} else {
		return false, User{}, fmt.Errorf("User validation failed")
	}
}

// Validate a user's current credentials, then change their password if they could be validated.
//
// Input:
//   - username, password string: User credentials. See ValidateUserCred.
//   - newPassword string: New password to set IF the above credentials can be validated.
// Output:
//   - error: Any error that occurs when changing user password, including: failure to validate,
//   user doesn't exist, failure to hash password
func ChangeUserPassword(username, password string, newPassword string) error {
	valid, _, err := ValidateUserCred(username, password)
	if !valid {
		if err != nil {
			return err
		} else {
			return fmt.Errorf("User validation failed")
		}
	}
	// If the user didn't exist, validation would have failed
	user, err := findUserEntryByUsername(username)
	if err != nil {
		return err
	}
	salt, err := genSalt()
	if err != nil {
		return err
	}
	pwdHash, err := slowHash([]byte(newPassword), salt, user.HashFunc)
	if err != nil {
		return err
	}
	user.PasswordHash = pwdHash
	user.Salt = stringEncode(salt)
	updateUser(&user)
	return nil
}

// Change user permissions for a given user.
// This action is generally initiated by an admin or the application server, and not a user; as a result,
// no password is required for the user.
//
// Input:
//   - username string: Username to alter.
//   - newPermissions map[string]bool: Full list of new permissions. Overwrites any permissions with the same name.
// Output:
//   - error: Any error that occurs while changing permissions, including: user does not exist, failure to marshal
//   permissions
func ChangeUserPermissions(username string, newPermissions map[string]bool) error {
	// Get userEntry. We need the entry ID to update the DB
	entry, err := findUserEntryByUsername(username)
	if err != nil {
		return err
	}
	// Get the corresponding User. This unmarshals the permissions string in entry so we can change permissions without a full overwrite.
	user := entry.toUser()
	// Write every key-value pair from newPermissions to old permissions.
	for name, value := range newPermissions {
		// It's possible user had no permissions before.
		if user.Permissions == nil {
			user.Permissions = make(map[string]bool)
		}
		user.Permissions[name] = value
	}
	// Re-stringify permissions and update the entry in the database.
	bytePerms, err := json.Marshal(user.Permissions)
	if err != nil {
		return err
	}
	entry.Permissions = string(bytePerms)
	updateUser(&entry)
	return nil
}
