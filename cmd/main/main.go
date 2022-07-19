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
		panic(err)
	}

	database.CloseDB()
}
