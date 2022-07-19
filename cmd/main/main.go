package main

import (
	"auth/pkg/database"
	"fmt"
)

func main() {
	/*
		pwd := []byte("passw0rd")
		start := time.Now()
		hash, salt, err := database.slowHash(pwd, "sha512")
		end := time.Now()
		if err != nil {
			panic(err)
		}
		fmt.Printf("HASH: %s\nSALT: %s\n", hash, salt)
		dur := end.Sub(start)
		fmt.Printf("Duration (ms): %d\n", dur.Milliseconds())
	*/
	// Try to validate fake user
	database.OpenDB()
	valid, err := database.ValidateUserCred("jnichols2719@protonmail.com", "jnichols", "password")
	if valid || err == nil {
		fmt.Println("This should have thrown an error...")
	}

	// Register me
	if err := database.RegisterUser("jnichols2719@protonmail.com", "jnichols", "password", nil); err == nil {
	} else {
		panic(err)
	}

	if valid, err := database.ValidateUserCred("jnichols2719@protonmail.com", "jnichols", "password"); err == nil {
		fmt.Printf("Validated: %t\n", valid)
	} else {
		panic(err)
	}

	if valid, err := database.ValidateUserCred("jnichols2719@protonmail.com", "jnichols", "wrongpassword"); err == nil {
		fmt.Printf("Validated: %t\n", valid)
	} else {
		panic(err)
	}

	// Try re-registering (should fail)
	// Register me
	if err := database.RegisterUser("jnichols2719@protonmail.com", "jnichols", "password", nil); err == nil {
		fmt.Println("Re-registration succeeded. It should not.")
	} else {
	}

	// Given that fails, let's change passwords. Wrong username?
	if err := database.ChangeUserPassword("jnichols2719@protonmail.com", "jnichols2719", "password", "newPassword"); err == nil {
		fmt.Println("Password change with wrong UN succeeded")
	}
	// Wrong old password
	if err := database.ChangeUserPassword("jnichols2719@protonmail.com", "jnichols", "wrongpassword", "newPassword"); err == nil {
		fmt.Println("Password change with wrong old password succeeded")
	}
	// Correct everything
	if err := database.ChangeUserPassword("jnichols2719@protonmail.com", "jnichols", "password", "newPassword"); err != nil {
		panic(err)
	}

	if valid, err := database.ValidateUserCred("jnichols2719@protonmail.com", "jnichols", "newPassword"); err == nil {
		fmt.Printf("Validated: %t\n", valid)
	} else {
		panic(err)
	}

	database.CloseDB()
}
