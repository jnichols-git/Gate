package authjwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// JWTHeader: `auth` treats these as constants.
type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// JWTBody: Token claims. All are registered except access, which is private.
type JWTBody struct {
	Issuer  string                 `json:"iss"`
	ForUser string                 `json:"sub"`
	Access  map[string]interface{} `json:"access"`
	Created int64                  `json:"iat"`
	Expires int64                  `json:"exp"`
}

// JSON Web Token structure combining the above.
type JSONWebToken struct {
	Header JWTHeader
	Body   JWTBody
}

// Create a new JWT based on a user email and access tag
func NewJWT(user string, access map[string]interface{}) *JSONWebToken {
	return &JSONWebToken{
		JWTHeader{
			Algorithm: "sha256",
			Type:      "jwt",
		},
		JWTBody{
			Issuer:  "auth",
			ForUser: user,
			Access:  access,
			Created: time.Now().Unix(),
			Expires: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
}

// Export a JSONWebToken using a given secret
func Export(t *JSONWebToken, secret []byte) (string, error) {
	h := hmac.New(sha256.New, secret)
	// Marshal and encode the JWT header/body separately
	head, _ := json.Marshal(t.Header)
	headStr := base64.RawURLEncoding.EncodeToString(head)
	body, err := json.Marshal(t.Body)
	if err != nil {
		return "", err
	}
	bodyStr := base64.RawURLEncoding.EncodeToString(body)
	// Write head.body to the hashing algorithm
	h.Write([]byte(
		headStr + "." + bodyStr,
	))
	// Get the signature from the hash
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	// Concatenate head.body.signature
	return headStr + "." + bodyStr + "." + signature, nil
}

// Verify that a token string is unaltered, unexpired, and signed with the given secret
func Verify(token string, secret []byte) (*JSONWebToken, bool, error) {
	items := strings.Split(token, ".")
	// Unmarshal and decode the JWT
	jwt := &JSONWebToken{}
	head, err := base64.RawURLEncoding.DecodeString(items[0])
	if err != nil {
		return nil, false, err
	}
	err = json.Unmarshal([]byte(head), &(jwt.Header))
	if err != nil {
		return nil, false, err
	}
	body, err := base64.RawURLEncoding.DecodeString(items[1])
	if err != nil {
		return nil, false, err
	}
	err = json.Unmarshal([]byte(body), &(jwt.Body))
	if err != nil {
		return nil, false, err
	}
	// Re-export the resulting jwt; should result in the exact same output
	expected, _ := Export(jwt, secret)
	// Return verification eval result and new token
	expired := jwt.Body.Expires < time.Now().Unix()
	return jwt, token == expected && !expired, nil
}
