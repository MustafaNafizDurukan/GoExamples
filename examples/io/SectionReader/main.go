package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fmt.Println("Read")
	Read()
	fmt.Println("------------------------------------------")
	fmt.Println("Seek")
	Seek()
	fmt.Println("------------------------------------------")
	fmt.Println("ReadAt")
	ReadAt()
	fmt.Println("------------------------------------------")
	fmt.Println("Size")
	Size()
}

func Read() {
	file, err := os.Open("myfile.txt")
	if err != nil {
		log.Fatal(err)
	}

	reader := io.NewSectionReader(file, 5, 100)

	firstChunk := make([]byte, 10)
	_, err = reader.Read(firstChunk)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(firstChunk))

	secondChunk := make([]byte, 15)
	_, err = reader.Read(secondChunk)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(secondChunk))

	lastChunk, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(lastChunk))
}

func Seek() {
	file, err := os.Open("myfile.txt")
	if err != nil {
		log.Fatal(err)
	}

	reader := io.NewSectionReader(file, 5, 100)

	// Move the offset 20 bytes from the start of the section
	_, err = reader.Seek(20, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	chunk := make([]byte, 15)
	_, err = reader.Read(chunk)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(chunk))
}

func ReadAt() {
	file, err := os.Open("myfile.txt")
	if err != nil {
		log.Fatal(err)
	}

	reader := io.NewSectionReader(file, 5, 100)

	chunk := make([]byte, 10)
	_, err = reader.ReadAt(chunk, 20)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(chunk))

	// Read the first 10 bytes of the section again
	_, err = reader.Read(chunk)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(chunk))
}

func Size() {
	file, err := os.Open("myfile.txt")
	if err != nil {
		log.Fatal(err)
	}

	reader := io.NewSectionReader(file, 5, 100)
	fmt.Println(reader.Size())
}
