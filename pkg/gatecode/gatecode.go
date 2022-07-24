// Package gatecode implements short-term code authentication, as opposed to gatekey's more long-term access.
//
// Generally speaking, gatecode should be used in tandem with a separate package capable of remote communication with users.
// gate provides [pkg/authmail] for this purpose. A regular control flow with gatecode may look like:
//   email := getUserEmailSomehow()
//   code := gatecode.NewGateCode(email)
//   sendEmailToUserSomehow(email, code)
//   ...
//   recEmail, recCode := receiveInputFromUserSomehow()
//   valid := gatecode.ValidateGateCode(recEmail, recCode)
// gatecode does not currently protect against abandoned codes, so there's a risk of filling up memory.
package gatecode

import (
	"math/rand"
	"time"
)

// An gateCode stores the email and code needed to validate, in addition to an expiration time.
type gateCode struct {
	Email   string
	Code    string
	Created time.Time
	Expires time.Time
}

// activeCodes maps emails to their respective codes.
var activeCodes map[string]*gateCode = make(map[string]*gateCode)

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
func NewGateCode(email string) string {
	now := time.Now()
	newCode := &gateCode{
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
func ValidateGateCode(email, code string) bool {
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
