package authcred

import (
	"fmt"

	"gorm.io/gorm"
)

var db *gorm.DB

// Add a userEntry to the database.
//
// Input:
//   - in *userEntry: User entry to add.
// Output:
//   - error: Returned if the database is closed.
func addUser(in *userEntry) error {
	if db == nil {
		return fmt.Errorf("addUser failed; database not open")
	}
	db.Create(in)
	return nil
}

// Update a userEntry in the database.
//
// Input:
//   - in *userEntry: User entry to alter. GORM manages this operation using primary key, which is the
//   int id of the entry. Highly suggest you don't try changing that.
// Output:
//   - error: Returned if the database is closed.
func updateUser(in *userEntry) error {
	if db == nil {
		return fmt.Errorf("updateUser failed; database not open")
	}
	db.Save(in)
	return nil
}

// Find a userEntry in the database by its email.
//
// Input:
//   - find string: Email to find.
// Output:
//   - out userEntry: The resulting userEntry.
//   - err error: Returned if the database is closed.
func findUserEntryByEmail(find string) (out userEntry, err error) {
	if db == nil {
		err = fmt.Errorf("findUser failed; database not open")
	} else {
		db.Where("email = ?", find).First(&out)
	}
	return
}

// Find a userEntry in the database by its username.
//
// Input:
//   - find string: Username to find.
// Output:
//   - out userEntry: The resulting userEntry.
//   - err error: Returned if the database is closed.
func findUserEntryByUsername(find string) (out userEntry, err error) {
	if db == nil {
		err = fmt.Errorf("findUser failed; database not open")
	} else {
		db.Where("Username = ?", find).First(&out)
	}
	return
}
