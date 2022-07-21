package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func OpenDB(path string) error {
	var err error
	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&UserEntry{})
	if err != nil {
		return err
	}

	return nil
}

func addUser(in *UserEntry) error {
	if db == nil {
		return fmt.Errorf("addUser failed; database not open")
	}
	db.Create(in)
	return nil
}

func updateUser(in *UserEntry) error {
	if db == nil {
		return fmt.Errorf("updateUser failed; database not open")
	}
	db.Save(in)
	return nil
}

func findUserByEmail(find string) (out UserEntry, err error) {
	if db == nil {
		err = fmt.Errorf("findUser failed; database not open")
	} else {
		db.Where("email = ?", find).First(&out)
	}
	return
}

func findUserByUsername(find string) (out UserEntry, err error) {
	if db == nil {
		err = fmt.Errorf("findUser failed; database not open")
	} else {
		db.Where("Username = ?", find).First(&out)
	}
	return
}
