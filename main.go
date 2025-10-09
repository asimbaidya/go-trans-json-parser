package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %q is not set or empty", key))
	}
	return value
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <word>")
		return
	}

	transPath := mustGetEnv("TRANS_DIR")

	inputWord := strings.ToLower(os.Args[1])
	filePath := filepath.Join(transPath, inputWord+".json")

	rawBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	term, err := parseWordJSON(inputWord, rawBytes)
	if err != nil {
		fmt.Println(err)
	}

	// print as json
	jsonBytes, err := json.MarshalIndent(term, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
