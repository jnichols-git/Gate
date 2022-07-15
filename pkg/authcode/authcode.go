package authcode

import (
	"fmt"
	"math/rand"
	"time"
)

type authorizationCode struct {
	ForUser string
	Code    string
	Created time.Time
	Expires time.Time
	Expired bool
}

var activeCodes map[string]*authorizationCode = make(map[string]*authorizationCode)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genCode(ct int) string {
	seq := make([]rune, ct)
	for i := range seq {
		seq[i] = letters[rand.Intn(len(letters))]
	}
	return string(seq)
}

// Create a new authorization code for a given user
func NewAuthCode(forUser string) *authorizationCode {
	now := time.Now()
	newCode := &authorizationCode{
		ForUser: forUser,
		Code:    genCode(6),
		Created: now,
		Expires: now.Add(time.Minute),
		Expired: false,
	}
	activeCodes[forUser] = newCode
	return newCode
}

// Validate an authorization code
func ValidateAuthCode(forUser, code string) bool {
	storedCode, ok := activeCodes[forUser]
	if !ok {
		return false
	}
	if storedCode.Code == code && storedCode.Expires.After(time.Now()) {
		delete(activeCodes, forUser)
		return true
	} else {
		fmt.Printf("Invalid: code equality %t, expired %t\n", storedCode.Code == code, !storedCode.Expires.After(time.Now()))
		delete(activeCodes, forUser)
		return false
	}
}
