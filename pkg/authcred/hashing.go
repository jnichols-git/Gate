package authcred

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"hash"
	"time"

	"github.com/pkg/errors"
)

// stringEncode shorthands for base64.URLEncoding.EncodeToString
func stringEncode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// stringDecode shorthands base64.URLEncoding.DecodeString
func stringDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// genSalt returns a 64-bit random byte string
func genSalt() ([]byte, error) {
	salt := make([]byte, 64)
	_, err := rand.Read(salt[:])
	if err != nil {
		return nil, err
	}
	return salt, nil
}

const hashRounds int = 65536
const minHashMS int = 20
const maxHashMS int = 0

var hfs map[string]func() hash.Hash = map[string]func() hash.Hash{
	"sha512": sha512.New,
}

// VERY SLOWLY hash pwd using a named hashFunc supported by auth.database, see hfs above
// Returns hash, salt, error (if occurs)
// DO NOT call any other hash function.
func slowHash(pwd, salt []byte, hashFunc string) (string, error) {
	// Get hash func. Return error if not supported
	f, ok := hfs[hashFunc]
	if !ok {
		return "", errors.Errorf("Hash function %s not supported", hashFunc)
	}
	hf := f()
	pwd_full := append(pwd, salt...)
	pwd_hash := make([]byte, base64.URLEncoding.EncodedLen(hf.Size()))
	// Initial hash into pwd_hash
	if _, err := hf.Write(pwd_full); err != nil {
		return "", err
	} else {
		pwd_hash = hf.Sum(nil)
	}
	// Repeatedly hash for hashRounds. Track time to complete.
	start := time.Now()
	for i := 0; i < hashRounds; i++ {
		if _, err := hf.Write(pwd_hash); err != nil {
			return "", err
		} else {
			pwd_hash = hf.Sum(nil)
		}
	}
	dur := time.Now().Sub(start).Milliseconds()
	if dur < int64(minHashMS) {
		return "", errors.Errorf("Hash function with %d rounds took %d ms (min %d). You should increase the number of rounds.", hashRounds, dur, minHashMS)
	}
	if dur > int64(maxHashMS) && maxHashMS > 0 {
		return "", errors.Errorf("Hash function with %d rounds took %d ms (max %d). This is not as bad as being too fast, but may increase latency.", hashRounds, dur, maxHashMS)
	}
	// Encode using base64URL
	encodedHash := stringEncode(pwd_hash)
	return encodedHash, nil
}
