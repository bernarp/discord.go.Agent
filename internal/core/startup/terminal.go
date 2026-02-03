package startup

import (
	"bufio"
	"fmt"
	"os"
)

func promptForToken() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--------------------------------------------------")
		fmt.Println("DISCORD BOT AGENT - STARTUP")
		fmt.Println("Discord Bot Token not found or invalid.")
		fmt.Print("Please enter a valid BOT_TOKEN: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("terminal input: %w", err)
		}

		token := CleanToken(input)

		if !ValidateToken(token) {
			fmt.Println("\n[ERROR] The token format is invalid.")
			fmt.Println("Example: MzE4... . ... . ...")
			continue
		}

		return token, nil
	}
}

func saveEncryptedToken(
	token string,
	key []byte,
) error {
	encrypted, err := encrypt([]byte(token), key)
	if err != nil {
		return fmt.Errorf("encryption: %w", err)
	}

	if err := os.WriteFile(TokenFileName, encrypted, 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	fmt.Println("✓ Token successfully encrypted and saved to", TokenFileName)
	fmt.Println("--------------------------------------------------")
	return nil
}
