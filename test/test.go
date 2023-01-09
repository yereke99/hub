package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
)

func encrypt(plain_text, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plain_text, nil), nil
}

func decrypt(cipher_text, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if len(cipher_text) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, cipher_text := cipher_text[:nonceSize], cipher_text[nonceSize:]

	return gcm.Open(nil, nonce, cipher_text, nil)
}

func main() {
	text := []byte("Hello, my name is Yerek.")

	key := []byte("the-key-has-to-be-32-bytes-long!")

	ciphertext, err := encrypt(text, key)

	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err.Error())
	}

	fmt.Printf("%s ==> %x\n", text, ciphertext)

	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err.Error())
	}

	fmt.Printf("%x ==> %s\n", ciphertext, plaintext)

}
