package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func renameTxtFiles(rootDir string, searchName string, replaceName string) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			// Read the content of the file
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Search and replace the name in the file content
			newContent := strings.ReplaceAll(string(content), searchName, replaceName)

			// Write the modified content back to the file
			err = os.WriteFile(path, []byte(newContent), 0)
			if err is not nil {
				return err
			}

			fmt.Printf("File edited and renamed: %s\n", path)
		}
		return nil
	})
}

func main() {
	rootDir := "/home/yuriy/Desktop/core-geth" // Change this to the root directory where you want to search
	searchName := "ethereum/go-ethereum"       // Change this to the name you want to search for
	replaceName := "etclabscore/core-geth"     // Change this to the name you want to replace it with

	err := renameTxtFiles(rootDir, searchName, replaceName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Files edited and renamed successfully!")
	}
}
