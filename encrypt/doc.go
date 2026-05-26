// Package encrypt provides AES-256-GCM authenticated encryption and decryption.
//
// The encrypted output is encoded as hex then base64, with the nonce prepended
// to the ciphertext. This package uses crypto/rand for nonce generation,
// ensuring each encryption produces unique output even for identical plaintext.
//
// Key requirements: AES-256 requires a 32-byte key. AES-128 (16-byte) and
// AES-192 (24-byte) keys are also accepted by the underlying cipher.
package encrypt
