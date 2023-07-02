package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fmt.Println("Write")
	Write()
	fmt.Println("------------------------------------------")
	fmt.Println("WriteAt")
	WriteAt()
	fmt.Println("------------------------------------------")
	fmt.Println("Seek")
	Seek()
}

func Write() {
	file, err := os.OpenFile("myfile.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	file.Write([]byte("Hi everyone! This is an example of using OffsetWriter in Go"))

	writer := io.NewOffsetWriter(file, 5)

	_, err = writer.Write([]byte("Hello, "))
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.Write([]byte("world!"))
	if err != nil {
		log.Fatal(err)
	}
}

func WriteAt() {
	file, err := os.OpenFile("myfile.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	file.Write([]byte("Hi everyone! This is an example of using OffsetWriter in Go"))

	writer := io.NewOffsetWriter(file, 5)

	_, err = writer.WriteAt([]byte("Hello, "), 10)
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.Write([]byte("world!"))
	if err != nil {
		log.Fatal(err)
	}
}

func Seek() {
	file, err := os.OpenFile("myfile.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	file.Write([]byte("Hi everyone! This is an example of using OffsetWriter in Go"))

	writer := io.NewOffsetWriter(file, 5)

	_, err = writer.Seek(10, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.Write([]byte("Hello, world!"))
	if err != nil {
		log.Fatal(err)
	}
}
