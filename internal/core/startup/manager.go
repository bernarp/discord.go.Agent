package startup

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetBotToken() (string, error) {
	for _, file := range EnvFiles {
		_ = godotenv.Load(file)
	}

	if token := CleanToken(os.Getenv(EnvTokenKey)); token != "" {
		if ValidateToken(token) {
			return token, nil
		}
		fmt.Printf("(!) Token in environment failed validation (Length: %d)\n", len(token))
	}

	key, err := deriveKey()
	if err != nil {
		return "", fmt.Errorf("startup key: %w", err)
	}

	if _, err := os.Stat(TokenFileName); err == nil {
		if token, err := tryLoadFromFile(key); err == nil {
			return token, nil
		}
	}

	token, err := promptForToken()
	if err != nil {
		return "", err
	}

	if err := saveEncryptedToken(token, key); err != nil {
		return "", err
	}

	return token, nil
}

func tryLoadFromFile(key []byte) (string, error) {
	data, err := os.ReadFile(TokenFileName)
	if err != nil {
		return "", err
	}

	decrypted, err := decrypt(data, key)
	if err != nil {
		fmt.Println("(!) Saved token decryption failed (hardware changed or file corrupted).")
		return "", err
	}

	token := CleanToken(string(decrypted))
	if !ValidateToken(token) {
		fmt.Println("(!) Decrypted token failed validation.")
		return "", fmt.Errorf("invalid file token")
	}

	return token, nil
}
