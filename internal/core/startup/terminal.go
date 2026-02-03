package startup

import (
	"bufio"
	"fmt"
	"os"
)

func promptGeneric(
	label, example string,
	optional bool,
	validator func(string) bool,
) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		optText := ""
		if optional {
			optText = " (optional, press Enter to skip)"
		}
		fmt.Printf("Please enter %s%s: ", label, optText)

		input, _ := reader.ReadString('\n')
		val := CleanInput(input)

		if val == "" && optional {
			return ""
		}

		if validator(val) {
			return val
		}

		fmt.Printf("[ERROR] Invalid %s format. Example: %s\n", label, example)
	}
}

func saveEncrypted(
	filename, value string,
	key []byte,
) error {
	if value == "" {
		return nil
	}
	encrypted, err := encrypt([]byte(value), key)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, encrypted, 0600)
}
