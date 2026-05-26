// Package hash provides password hashing and verification using industry-standard
// algorithms: Argon2id, bcrypt, and SHA-256 (with optional HMAC key).
//
// For password storage, prefer Argon2id (NewArgon2/CompareArgon2) or bcrypt
// (CreateBcrypt/CompareBcrypt). SHA-256 is provided for non-password hashing
// use cases such as data integrity checks.
//
// All hashing functions use crypto/rand for salt generation and constant-time
// comparison to prevent timing attacks.
package hash
