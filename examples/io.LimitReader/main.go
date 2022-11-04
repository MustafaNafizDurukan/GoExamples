package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func main() {
	sr := strings.NewReader("io.limitReader example")
	r := io.LimitReader(sr, 8)

	br := bufio.NewReader(r)

	buf := make([]byte, 30)
	n, err := br.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Whole byte size: %d. Read byte size: %d \n", len(buf), n)
	fmt.Printf("Read byte: %s \n", buf)
}
