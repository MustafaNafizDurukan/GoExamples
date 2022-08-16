package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mustafanafizdurukan/goexamples/pkg/constants"
)

func main() {
	res, err := http.Get("http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-all-5gram-20120701-0.gz")
	if err != nil {
		panic(err)
	}

	local, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer local.Close()

	dec, err := gzip.NewReader(res.Body)
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(local, io.TeeReader(dec, &progress{})); err != nil {
		panic(err)
	}
}

type progress struct {
	total uint64
}

func (p *progress) Write(b []byte) (int, error) {
	p.total += uint64(len(b))
	progress := float64(p.total) / (constants.MB)

	fmt.Printf("Downloading %f MB... \n", progress)
	return len(b), nil
}
