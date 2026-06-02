package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

func Decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

// DecryptAES256 decrypt data yang dienkripsi dengan:
// - AES-256-CBC
// - key di-hash SHA256
// - IV = 16 byte pertama dari data
func DecryptAES256(encryptedString string, keyString string) (string, error) {
	// SHA256 key → 32 byte
	key := sha256.Sum256([]byte(keyString))

	// Base64 decode
	combined, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	if len(combined) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// Ambil IV & ciphertext
	iv := combined[:aes.BlockSize]         // 16 byte
	ciphertext := combined[aes.BlockSize:] // sisanya

	// Buat AES block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt (in-place)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Unpadding PKCS7
	plaintext, err := pkcs7Unpad(ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}

	paddingLen := int(data[length-1])
	if paddingLen > length {
		return nil, errors.New("invalid padding")
	}

	return data[:length-paddingLen], nil
}
