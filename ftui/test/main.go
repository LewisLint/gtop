package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path")
		return
	}

	filePath := os.Args[1]

	// 1. Open the file
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 2. Read the first 512 bytes into a buffer
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}

	// 3. Detect content type using the actual bytes read
	fileType := http.DetectContentType(buffer[:n])
	fmt.Printf("Filetype: %s\n", fileType)
}
