# SectionReader

Hello,

In this article we are gonna talk about SectionReader in io standart package. We are gonna talk about What SectionReader actually is? Why does it exist? and How can we use it?

Lets dive in.

**myfile.txt**

```go
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut auctor justo eu ipsum dapibus, sit amet sagittis erat interdum.
```

# What?

```go
// SectionReader is a Reader that reads from a section of a file.
type SectionReader struct {
	r     ReaderAt
	base  int64
	off   int64
	limit int64
	err   error
}

// NewSectionReader returns a SectionReader that reads from r starting at
// offset off and stops with EOF after n bytes.
func NewSectionReader(r ReaderAt, off, n int64) *SectionReader {
	return &SectionReader{r, off, off, off + n, nil}
}

// Read implements the io.Reader interface.
func (s *SectionReader) Read(p []byte) (n int, err error) {
	if s.err != nil {
		return 0, s.err
	}
	if s.off >= s.limit {
		return 0, io.EOF
	}
	if max := s.limit - s.off; int64(len(p)) > max {
		p = p[:max]
	}
	n, s.err = s.r.ReadAt(p, s.off)
	s.off += int64(n)
	return
}

// Seek implements the io.Seeker interface.
func (s *SectionReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errors.New("Seek: invalid whence")
	case io.SeekStart:
		offset += s.base
	case io.SeekCurrent:
		offset += s.off
	case io.SeekEnd:
		offset += s.limit
	}
	if offset < s.base {
		return 0, errors.New("Seek: negative position")
	}
	s.off = offset
	return offset - s.base, nil
}

// ReadAt implements the io.ReaderAt interface.
func (s *SectionReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off >= s.limit-s.base {
		return 0, io.EOF
	}
	off += s.base
	if max := s.limit - off; int64(len(p)) > max {
		p = p[:max]
	}
	return s.r.ReadAt(p, off)
}

```

Above code, you can see direct code of SectionReader implementation. As you can see we have Read, Seek, ReadAt functions but these are not just random functions. These are implementation of **`io.Reader`**, **`io.Seeker`**, and **`io.ReaderAt`** interfaces**.**

If we simply list the features of the above code into items:

- It implements the **`io.Reader`**, **`io.Seeker`**, and **`io.ReaderAt`** interfaces, so it can be used as a drop-in replacement for other types that implement these interfaces.
- **`NewSectionReader()`** function, takes a **`ReaderAt`** instance and the starting and ending positions of the section of the file to read. It returns a **`SectionReader`** instance that reads from the specified section of the file.
- **`Read()`** method, reads data from the section of the file and returns the number of bytes read and any error that occurred.
- `Seek()` method, sets the current position in the file to a specified offset and returns the new position in the file and any error that occurred. This method allows the `SectionReader` to implement the `io.Seeker` interface.
- `ReadAt()` method, reads data from the specified position in the file. This method allows the `SectionReader` to implement the `io.ReaderAt` interface.

As you can see above, there are 2 read functions which names are `Read` and `ReadAt` but they have been created for same purpose which is reading n byte at specific offset. Also there is Seek function allows us to continue reading from the nth byte.

# Why?

The **`SectionReader`** type in Go provides a way to read from a section of an underlying **`ReaderAt`** object, allowing you to read a portion of a file or other data stream without reading the entire thing into memory. This can be useful when you want to save memory and improve performance in your Go programs.

To read data from a **`SectionReader`** object, you can use the **`Read`** method provided by the **`SectionReader`** type. The **`Read`** method reads data from the underlying **`ReaderAt`** object, starting at the offset specified when the **`SectionReader`** was created and ending after the number of bytes specified when the **`SectionReader`** was created. This allows you to read a section of a file or other data stream without reading the entire thing into memory.

If the **`Read`** method reaches the end of the section specified when the **`SectionReader`** was created, it returns an **`EOF`** error, indicating that there is no more data to read. This allows you to easily detect when the end of the section has been reached and stop reading.

Overall, the **`Read`** method of the **`SectionReader`** type is a convenient and efficient way to read a section of a file or other data stream in Go, allowing you to save memory and improve performance in your Go programs.

# How?

```go
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
```

1. We open file which we want to read.
2. Create SectionReader to specify where will we begin to read and how many byte we will read.
    - In this example we will begin 10th byte and we will 20 byte.
3. We give SectionReader to io.ReadAll and we call io.ReadAll. In this way, it calls the Read function of the SectionReader. 
    
    ```go
    func ReadAll(r Reader) ([]byte, error) {
    	b := make([]byte, 0, 512)
    	for {
    		if len(b) == cap(b) {
    			// Add more capacity (let append pick how much).
    			b = append(b, 0)[:len(b)]
    		}
    		n, err := r.Read(b[len(b):cap(b)])
    		b = b[:len(b)+n]
    		if err != nil {
    			if err == EOF {
    				err = nil
    			}
    			return b, err
    		}
    	}
    }
    ```
    
    Now we are at io.ReadAll function. As you can see here this function creates byte slice which lentgh is 512 and calls Read function of given Reader which is SectionReader’s Read function.
    
    ```go
    func (s *SectionReader) Read(p []byte) (n int, err error) {
    	if s.off >= s.limit {
    		return 0, EOF
    	}
    	if max := s.limit - s.off; int64(len(p)) > max {
    		p = p[0:max]
    	}
    	n, err = s.r.ReadAt(p, s.off)
    	s.off += int64(n)
    	return
    }
    ```
    
    1. Now we are at SectionReader’s Read function. As you can see we will byte slice which has been created by io.ReadAll function and we will reslice it to given number of byte before which was 20. 
    2. We resliced the slice back to size 20 and we called s.r.ReadAt which is os.File.ReadAt function.
4. Now we are printing `data` (slice)
    
    ```go
    20
    m dolor sit amet, co
    ```
    
5. We seeked for beginning of offset we gave so we are 10th byte again right now.
6. When we retry to read program will return same output as before because as i said we seeked beginning of the offset and we read it again.
7. Last but not least, output of io.ReadAll function is

```go
m dolor sit amet, co
```

Thanks for your time.

Happy coding!