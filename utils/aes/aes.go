package aes

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/hex"
	"io"

	"github.com/pkg/errors"
)

// Encrypt 加密
func Encrypt(key, plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "utils: AesEncrypt failed")
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return "", errors.Wrap(err, "utils, AesEncrypt failed")
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(plaintext))
	return hex.EncodeToString(ciphertext), nil

}

// Decrypt 解密
func Decrypt(key, d string) (string, error) {
	ciphertext, err := hex.DecodeString(d)
	if err != nil {
		return "", errors.Wrap(err, "utils: AesDecrypt failed")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "utils: AesDecrypt failed")
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("utils: AesDecrypt ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
