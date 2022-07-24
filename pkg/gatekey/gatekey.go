package gatekey

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// jwtHeader: `auth` treats these as constants.
type GateKeyHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// jwtBody: Token claims. All are registered except access, which is private.
type GateKeyBody struct {
	Issuer      string          `json:"iss"`
	ForUser     string          `json:"sub"`
	Permissions map[string]bool `json:"permissions"`
	Created     int64           `json:"iat"`
	Expires     int64           `json:"exp"`
}

// Gate key structure.
// GateKeys are a valid JWT struct, and can be exported/validated as one.
type GateKey struct {
	Header GateKeyHeader
	Body   GateKeyBody
}

// Generate a new GateKey.
//
// Input:
//   - username string: Attaches this key to a specific user.
//   - permissions map[string]bool: Represents the user permissions that this key grants
//   - expiry time.Duration: A *duration* for how long the token should last.
// Output:
//   - *jwt
func NewGateKey(username string, permissions map[string]bool, expiry time.Duration) *GateKey {
	return &GateKey{
		GateKeyHeader{
			Algorithm: "sha256",
			Type:      "jwt",
		},
		GateKeyBody{
			Issuer:      "auth",
			ForUser:     username,
			Permissions: permissions,
			Created:     time.Now().Unix(),
			Expires:     time.Now().Add(expiry).Unix(),
		},
	}
}

// Export a GateKey.
// This and Verify are inverse operations; the same secret MUST be used in both for correct results.
//
// Input:
//   - key *GateKey: Key to export. Should be non-nil.
//   - secret []byte: Signing secret.
// Output:
//   - string: Exported key
func Export(key *GateKey, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	// Marshal and encode the JWT header/body separately
	head, _ := json.Marshal(key.Header)
	headStr := base64.RawURLEncoding.EncodeToString(head)
	body, _ := json.Marshal(key.Body)
	bodyStr := base64.RawURLEncoding.EncodeToString(body)
	// Write head.body to the hashing algorithm
	h.Write([]byte(
		headStr + "." + bodyStr,
	))
	// Get the signature from the hash
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	// Concatenate head.body.signature
	return headStr + "." + bodyStr + "." + signature
}

// Verify an exported GateKey.
// This and Export are inverse operations; the same secret MUST be used in both for correct results.
// Notably, this returns a bool representing validity AND an error. The error does not represent the state
// of the token, but any error that occured during validation. It is recommended that when checking if the
// result of Verify, if the bool output is false, that the error also be checked for a reason; if the error
// is nil, then the token has a valid format but is either not signed correctly or expired.
//
// Input:
//   - token string: Exported GateKey.
//   - secret []byte: Signing secret.
// Output:
//   - *GateKey: Resulting GateKey. nil if verification failed.
//   - bool: Is token valid?
//   - error: Any error that occurs during verification.
func Verify(token string, secret []byte) (*GateKey, bool, error) {
	items := strings.Split(token, ".")
	// Unmarshal and decode the JWT
	key := &GateKey{}
	head, err := base64.RawURLEncoding.DecodeString(items[0])
	if err != nil {
		return nil, false, err
	}
	err = json.Unmarshal([]byte(head), &(key.Header))
	if err != nil {
		return nil, false, err
	}
	body, err := base64.RawURLEncoding.DecodeString(items[1])
	if err != nil {
		return nil, false, err
	}
	err = json.Unmarshal([]byte(body), &(key.Body))
	if err != nil {
		return nil, false, err
	}
	// Re-export the resulting key; should result in the exact same output
	expected := Export(key, secret)
	// Return verification eval result and new token
	expired := key.Body.Expires < time.Now().Unix()
	return key, token == expected && !expired, nil
}
