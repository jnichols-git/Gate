package authcode

import (
	"math/rand"
	"time"
)

type authorizationCode struct {
	Email   string
	Code    string
	Created time.Time
	Expires time.Time
}

var activeCodes map[string]*authorizationCode = make(map[string]*authorizationCode)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

/* Generates a ct-length authorization code using all-caps letters.
Input:
	ct int: Number of characters to generate
Output:
	string: string authorization code
*/
func genCode(ct int) string {
	seq := make([]rune, ct)
	for i := range seq {
		seq[i] = letters[rand.Intn(len(letters))]
	}
	return string(seq)
}

/* Generates a new authorization code, and stores it in memory for checking later.
Input:
	email string: Output code will be attached to this email. See authmail for how mail is sent.
Output:
	string: string authorization code.
*/
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

/* Validates a given authorization code against an email.
Input:
	email, code string: Both the email and code must match records.
Output:
	bool: Represents code validity. true if code correct and unexpired, false if email incorrect, code incorrect, or expired.
*/
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
