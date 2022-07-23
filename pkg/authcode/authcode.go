// Package authcode implements short-term code authentication.
//
// Generally speaking, authcode should be used in tandem with a separate package capable of remote communication with users.
// auth provides [pkg/authmail] for this purpose. A regular control flow with authcode may look like:
//   email := getUserEmailSomehow()
//   code := authcode.NewAuthCode(email)
//   sendEmailToUserSomehow(email, code)
//   ...
//   recEmail, recCode := receiveInputFromUserSomehow()
//   valid := authcode.ValidateAuthCode(recEmail, recCode)
// authcode does not currently protect against abandoned codes, so there's a risk of filling up memory.
package authcode

import (
	"math/rand"
	"time"
)

// An authorizationCode stores the email and code needed to validate, in addition to an expiration time.
type authorizationCode struct {
	Email   string
	Code    string
	Created time.Time
	Expires time.Time
}

// activeCodes maps emails to their respective codes.
var activeCodes map[string]*authorizationCode = make(map[string]*authorizationCode)

// letters contains a list of valid runes used in authorization codes.
var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Generates a ct-length authorization code using all-caps letters.
//
// Input:
//   - ct int: Number of characters to generate
// Output:
//   - string: string authorization code
func genCode(ct int) string {
	seq := make([]rune, ct)
	for i := range seq {
		seq[i] = letters[rand.Intn(len(letters))]
	}
	return string(seq)
}

// Generates a new authorization code, and stores it in memory for checking later.
//
// Input:
//   - email string: Output code will be attached to this email. See authmail for how mail is sent.
// Output:
//   - string: string authorization code.
func NewAuthCode(email string) string {
	now := time.Now()
	newCode := &authorizationCode{
		Email:   email,
		Code:    genCode(6),
		Created: now,
		Expires: now.Add(time.Minute * 5),
	}
	activeCodes[email] = newCode
	return newCode.Code
}

// Validates a given authorization code against an email.
//
// Input:
//   - email, code string: Both the email and code must match records.
// Output:
//   - bool: Represents code validity. true if code correct and unexpired, false if email incorrect, code incorrect, or expired.
func ValidateAuthCode(email, code string) bool {
	storedCode, ok := activeCodes[email]
	if !ok {
		return false
	}
	if storedCode.Code == code && storedCode.Expires.After(time.Now()) {
		delete(activeCodes, email)
		return true
	} else {
		delete(activeCodes, email)
		return false
	}
}
