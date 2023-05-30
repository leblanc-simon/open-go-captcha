package crypto

// https://earthly.dev/blog/cryptography-encryption-in-go/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"leblanc.io/open-go-captcha/config"
)

var secretKey []byte

func Initialize(c *config.Config) {
	hash := sha256.New()
	hash.Write([]byte(c.Captcha.SecretKey))
	bs := hash.Sum(nil)

	secretKey = bs
}

func Encrypt(value string) (string) {
	aesBlock, err := aes.NewCipher(secretKey)
	if err != nil {
	 fmt.Println(err)
	}
   
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
	 fmt.Println(err)
	}
	
	nonce := make([]byte, gcmInstance.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	cipheredText := gcmInstance.Seal(nonce, nonce, []byte(value), nil)

	return base64.URLEncoding.EncodeToString(cipheredText)
}

func Decrypt(value string) (string, error) {
	ciphered, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	
	aesBlock, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return "", err
	}

	return string(originalText), nil
}