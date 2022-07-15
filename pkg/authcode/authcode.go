package authcode

import (
	"fmt"
	"math/rand"
	"time"
)

type AuthCode struct {
	ForUser string
	Code    string
	Created time.Time
	Expires time.Time
	Expired bool
}

var activeCodes map[string]*AuthCode = make(map[string]*AuthCode)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genCode(ct int) string {
	seq := make([]rune, ct)
	for i := range seq {
		seq[i] = letters[rand.Intn(len(letters))]
	}
	return string(seq)
}

func NewAuthCode(forUser string) *AuthCode {
	now := time.Now()
	newCode := &AuthCode{
		ForUser: forUser,
		Code:    genCode(6),
		Created: now,
		Expires: now.Add(time.Minute),
		Expired: false,
	}
	activeCodes[forUser] = newCode
	return newCode
}

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
