package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type Person struct {
	Name    string
	Surname string
}

func (p *Person) Write(b []byte) (n int, err error) {
	fmt.Printf("Write method ran. Person's name is %s \n", p.Name)
	return len(b), nil
}

func main() {
	persons := []Person{
		{Name: "name1", Surname: "surname1"},
		{Name: "name2", Surname: "surname2"},
		{Name: "name3", Surname: "surname3"},
		{Name: "name4", Surname: "surname4"},
	}

	f, err := os.Create("1.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := bytes.NewBuffer(make([]byte, 0))

	w := io.MultiWriter(f, buf, os.Stdout)
	for i := range persons {
		w = io.MultiWriter(w, &persons[i])
	}

	enc := xml.NewEncoder(w)
	if err := enc.Encode(persons); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Buffer string:", buf.String())
}
