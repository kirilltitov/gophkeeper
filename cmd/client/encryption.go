package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
)

func getEncryptionKeyBytes(cmd *cli.Command, force bool) ([]byte, error) {
	if cmd.Bool(flagNoEncrypt) && !force {
		fmt.Fprintf(cmd.Root().Writer, "WARNING: You have disabled encryption key prompt, this might be unsecure\n")
		return nil, nil
	}

	encryptionKeyString, err := readPassword(cmd.Root().Writer, "Enter encryption key (NOT PASSWORD): ")

	if encryptionKeyString == "" {
		fmt.Fprintf(cmd.Root().Writer, "WARNING: You provided an empty encryption key, this might be unsecure\n")
		return nil, nil
	}

	result := sha256.Sum256([]byte(encryptionKeyString))

	return result[:], err
}

func encrypt(keyBytes []byte, input []byte) (string, error) {
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	encryptedBytes := gcm.Seal(nil, nonce, input, nil)
	encryptedBytes = append(encryptedBytes, nonce...)

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func decrypt(keyBytes []byte, text string) ([]byte, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := encryptedBytes[len(encryptedBytes)-gcm.NonceSize():]
	decryptedBytes, err := gcm.Open(nil, nonce, encryptedBytes[:len(encryptedBytes)-gcm.NonceSize()], nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not decrypt secret")
	}

	return decryptedBytes, nil
}
