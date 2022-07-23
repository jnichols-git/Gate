package authcred

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"hash"
	"time"

	"github.com/pkg/errors"
)

// stringEncode is shorthand for base64.URLEncoding.EncodeToString().
var stringEncode func([]byte) string = base64.RawURLEncoding.EncodeToString

// stringDecode is shorthand for base64.URLEncoding.DecodeString().
var stringDecode func(string) ([]byte, error) = base64.RawURLEncoding.DecodeString

// genSalt returns a 64-bit random byte string
func genSalt() ([]byte, error) {
	salt := make([]byte, 64)
	_, err := rand.Read(salt[:])
	if err != nil {
		return nil, err
	}
	return salt, nil
}

const hashRounds int = 131072 // 2^17
const minHashMS int = 20

// hfs is a map of string hash function names to a Go Hash object creator.
// See slowHash for an example of usage.
var hfs map[string]func() hash.Hash = map[string]func() hash.Hash{
	"sha512": sha512.New,
}

// Hashes a byte string VERY SLOWLY with a specific hashfunc supported by hfs.
// Any private value used in auth MUST be hashed through slowHash.
//
// Input:
//   - pwd, salt byte: Password (or other string) to hash, and the salt to hash it with.
//   - hashFunc string: Hash function to use. Must be a key in hfs.
// Output:
//   - string: Output hashed value
//   - error: Any error that occurs, including: unsupported hash function or too-efficient hashing.
//   Passwords will be put through hashRounds rounds, and if that process takes <20ms, that throws an error.
//   It is Very Bad if this ever happens.
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
		return "", errors.Errorf("Hash function with %d rounds took %d ms (min %d).", hashRounds, dur, minHashMS)
	}
	// Encode using base64URL
	encodedHash := stringEncode(pwd_hash)
	return encodedHash, nil
}
