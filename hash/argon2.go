package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"golang.org/x/crypto/argon2"
)

const (
	standardLengthValues = 6

	/*
		NOTE: values with '$' as delimiter
			1: hash algorithm name
			2: argon2 (v)ersion
			3: (m)emory, i(t)erations, (p)arallelism used for hash
			4: salt with encoded base64
			5: hash with encoded base64
	*/
	standardHashFormat = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
)

type params struct {
	salt        []byte
	iterations  uint32
	memory      uint32
	parallelism uint8
	keyLen      uint32
}

// hash password with argon2 standard hash format
func NewArgon2(password []byte) (string, error) {
	p := params{
		salt:        make([]byte, 16),
		iterations:  3,
		memory:      4 * 1024,
		parallelism: 1,
		keyLen:      32,
	}

	_, err := rand.Read(p.salt)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("error while read random number for salt for argon2, %s", err))
	}

	hash := argon2.IDKey(password, p.salt, p.iterations, p.memory, p.parallelism, p.keyLen)
	encodedHash, err := encodeHash(&p, hash)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeArgon2EncodeHashError, fmt.Sprintf("error while encode argon2 hash, %s", err.Error()))
	}

	return encodedHash, nil
}

// compare hash password with argon2 standard hash format
func CompareArgon2(password, encodedHash string) (bool, error) {
	// extract all parameters include salt and key length from encoded password hash
	p, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// generate other password with sam parameters
	otherHash := argon2.IDKey([]byte(password), p.salt, p.iterations, p.memory, p.parallelism, p.keyLen)

	// check that the contents of the hashed passwords are identical
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// encode argon2 hash to standard argon2 format hash
func encodeHash(p *params, hash []byte) (string, error) {
	b64Salt := base64.RawStdEncoding.EncodeToString(p.salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(standardHashFormat, argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

// decode argon2 format hash to value parameter
func decodeHash(encodedHash string) (p *params, hash []byte, err error) {
	values := strings.Split(encodedHash, "$")
	if lenValues := len(values); lenValues != standardLengthValues {
		return nil, nil, errors.NewWithCode(codes.CodeArgon2InvalidEncodedHash, fmt.Sprintf("invalid length of encoded hash, expected %v but get %v", standardLengthValues, lenValues))
	}

	// check compatible argon2 version
	version := 0
	if _, err := fmt.Sscanf(values[2], "v=%d", &version); err != nil {
		return nil, nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("error while get argon2 version on decode hash, %s", err))
	}
	if version != argon2.Version {
		return nil, nil, errors.NewWithCode(codes.CodeArgon2IncompatibleVersion, fmt.Sprintf("incompatible argon2 version, current argon2 version is %d used to compare with version %d", argon2.Version, version))
	}

	// mapping values for memory, iterations and parallelism
	p = &params{}
	if _, err := fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism); err != nil {
		return nil, nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("error while get values from memory, iterations and parallelism, %s", err))
	}

	// decode base64 salt
	p.salt, err = base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("error while decode salt from values, %s", err))
	}

	// decode base64 hash
	hash, err = base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("error while decode hash from values, %s", err))
	}
	p.keyLen = uint32(len(hash))

	return p, hash, nil
}
