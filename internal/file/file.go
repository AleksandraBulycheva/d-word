package file

import (
	"os"
)

// ReadFile reads the content of a file
func ReadFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []byte{}, nil
		}
		return nil, err
	}
	return data, nil
}

// WriteFile writes content to a file
func WriteFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}
