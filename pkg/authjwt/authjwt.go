package authjwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"
)

type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type JWTBody struct {
	Issuer  string `json:"iss"`
	ForUser string `json:"sub"`
	Access  string `json:"access"`
	Created int64  `json:"iat"`
	Expires int64  `json:"exp"`
}

type JSONWebToken struct {
	Header    JWTHeader
	Body      JWTBody
	Signature string
}

// Create a new JWT based on a user email and access tag
func NewJWT(user, access string) *JSONWebToken {
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
		"",
	}
}

// Create a JWT struct by unmarshaling a string token
func UnmarshalJWT(token string) *JSONWebToken {
	t := &JSONWebToken{}
	json.Unmarshal([]byte(token), t)
	return t
}

// Sign a JWT using a secret
func (t *JSONWebToken) Sign(secret []byte) {
	h := hmac.New(sha256.New, secret)
	head, _ := json.Marshal(t.Header)
	body, _ := json.Marshal(t.Body)
	h.Write([]byte(
		base64.RawURLEncoding.EncodeToString(head) + "." + base64.RawURLEncoding.EncodeToString(body),
	))
	t.Signature = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// Verify the signature on a given JWT
func VerifySignature(t *JSONWebToken, secret []byte) bool {
	h := hmac.New(sha256.New, secret)
	head, _ := json.Marshal(t.Header)
	body, _ := json.Marshal(t.Body)
	h.Write([]byte(
		base64.RawURLEncoding.EncodeToString(head) + "." + base64.RawURLEncoding.EncodeToString(body),
	))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return signature == t.Signature
}
