package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// Create a new io.Reader that reads from a file
	file, err := os.Open("myfile.txt")
	if err != nil {
		// Handle the error
	}

	// Create a new io.SectionReader that reads from the file starting at byte 10
	// and ending at byte 30
	reader := io.NewSectionReader(file, 10, 20)

	// Read the contents of the io.SectionReader using the Read method
	data, err := io.ReadAll(reader)
	if err != nil {
		// Handle the error
	}

	// Use the data read from the io.SectionReader
	fmt.Println(len(data))
	fmt.Println(string(data))

	// Seek to the beginning of the io.SectionReader using the Seek method
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		// Handle the error
	}

	// Read the contents of the io.SectionReader again using the ReadAt method
	data, err = io.ReadAll(reader)
	if err != nil {
		// Handle the error
	}

	// Use the data read from the io.SectionReader
	fmt.Println(string(data))
}
