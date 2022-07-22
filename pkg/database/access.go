package database

import (
	"fmt"

	"gorm.io/gorm"
)

var db *gorm.DB

func addUser(in *userEntry) error {
	if db == nil {
		return fmt.Errorf("addUser failed; database not open")
	}
	db.Create(in)
	return nil
}

func updateUser(in *userEntry) error {
	if db == nil {
		return fmt.Errorf("updateUser failed; database not open")
	}
	db.Save(in)
	return nil
}

func findUserByEmail(find string) (out userEntry, err error) {
	if db == nil {
		err = fmt.Errorf("findUser failed; database not open")
	} else {
		db.Where("email = ?", find).First(&out)
	}
	return
}

func findUserByUsername(find string) (out userEntry, err error) {
	if db == nil {
		err = fmt.Errorf("findUser failed; database not open")
	} else {
		db.Where("Username = ?", find).First(&out)
	}
	return
}
