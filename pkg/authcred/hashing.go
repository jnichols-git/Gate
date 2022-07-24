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

// Creates a length 64 random byte string.
//
// Output:
//   - []byte: Byte string
//   - error: Currently always nil. Error handling for read/write calls is currently a project-wide issue.
func genSalt() ([]byte, error) {
	salt := make([]byte, 64)
	bytesRead := 0
	// Make 3 attempts to generate salt value.
	for attempts := 0; attempts < 3 && bytesRead < 64; attempts++ {
		attBytesRead, err := rand.Read(salt[bytesRead:])
		bytesRead += attBytesRead
		// We do need to check the error here, but it may not be fatal (EOF). This behavior could be improved.
		// It is usually OS-related and asks the user to wait for the syscall to be available for random generation,
		// so if an error occurs and we didn't finish reading to salt, we'll wait for 100ms
		if err != nil && bytesRead < 64 {
			time.Sleep(time.Millisecond * 100)
		}
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
//   - error: Any error that occurs, including: unsupported hash function
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
	// Repeatedly hash for 2^17 rounds.
	for i := 0; i < 131072; i++ {
		if _, err := hf.Write(pwd_hash); err != nil {
			return "", err
		} else {
			pwd_hash = hf.Sum(nil)
		}
	}
	// Encode using base64URL
	encodedHash := stringEncode(pwd_hash)
	return encodedHash, nil
}
