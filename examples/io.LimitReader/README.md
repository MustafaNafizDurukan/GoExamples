# LimitReader

## Explaining

In this reading, we want to see how io.LimitReader works so i prepared an example for you. 

We will create strings.Reader and we will read this bytes up to a certain number of bytes so lets begin.

```go
sr := strings.NewReader("io.limitReader example")
```

First of all, we will create strings.Reader with given string.

```go
r := io.LimitReader(sr, 8)
```

then  we create io.Reader with the number we gave to LimitReader.

```go
br := bufio.NewReader(r)
```

Now we define bufio.Reader which implements and contains io.Reader interface.

```go
buf := make([]byte, 30)
n, err := br.Read(buf)
if err != nil {
	fmt.Println(err)
	return
}
```

finally we call bufio.Read function and send byte in it. When we call bufio.Reader, it will check some conditions after that it will call given reader which is io.LimitReader.

```go
// LimitReader returns a Reader that reads from r
// but stops with EOF after n bytes.
// The underlying implementation is a *LimitedReader.
func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }

// A LimitedReader reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns EOF when N <= 0 or when the underlying R returns EOF.
type LimitedReader struct {
	R Reader // underlying reader
	N int64  // max bytes remaining
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}
```

As you see above, when we call io.LimitReader function it actually creates io.LimitedReader which contains io.Reader and int that specifies max bytes to be read. 

When you call its Read function, it reslices given slice and calls Read function of strings package.

## Debugging

---

Before starting, i want to say that I put debug points on the functions below.

- (b *Reader) Read(p []byte) (n int, err error)  (bufio package)
- (l *LimitedReader) Read(p []byte) (n int, err error)  (io package)
- (r *Reader) Read(b []byte) (n int, err error)  (strings package)

lets start.

When i run the program first breakpoint that hit is `Read` function in bufio package.

```go
// Read reads data into p.
// It returns the number of bytes read into p.
// The bytes are taken from at most one Read on the underlying Reader,
// hence n may be less than len(p).
// To read exactly len(p) bytes, use io.ReadFull(b, p).
// If the underlying Reader can return a non-zero count with io.EOF,
// then this Read method can do so as well; see the [io.Reader] docs.
func (b *Reader) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		if b.Buffered() > 0 {
			return 0, nil
		}
		return 0, b.readErr()
	}
	if b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		if len(p) >= len(b.buf) {
			// Large read, empty buffer.
			// Read directly into p to avoid copy.
			n, b.err = b.rd.Read(p)
			if n < 0 {
				panic(errNegativeRead)
			}
			if n > 0 {
				b.lastByte = int(p[n-1])
				b.lastRuneSize = -1
			}
			return n, b.readErr()
		}
		// One read.
		// Do not use b.fill, which will loop.
		b.r = 0
		b.w = 0
		n, b.err = b.rd.Read(b.buf)
		if n < 0 {
			panic(errNegativeRead)
		}
		if n == 0 {
			return 0, b.readErr()
		}
		b.w += n
	}

	// copy as much as we can
	// Note: if the slice panics here, it is probably because
	// the underlying reader returned a bad count. See issue 49795.
	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	b.lastByte = int(b.buf[b.r-1])
	b.lastRuneSize = -1
	return n, nil
}
```

After conditions and settings variables, we call Read func of read in bufio.Read struct which is actually Read func of limitReader.

Run program again.

```go
func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}
```

Now we hit LimitReader.Read. As you can see, this function takes maximum number of bytes to be read and reslice the given slice to that number and after that calls actual Read function with new slice.

Rerun program.

 

```go
type Reader struct {
	s        string
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}

// Read implements the io.Reader interface.
func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}
```

Now finally we hit strings.Reader function. This function copies string in Reader starting with i to the given bytes.

Output

```go
Whole byte size: 30. Read byte size: 8 
Read byte: io.limit
```

If i call `n, err := br.Read(buf)` once more, i could not read anything because limitreader prevents me doing that. 

There is max bytes remaining field in LimitReader. After each read, this field is reduced by the number of bytes read so when this function is called a second time, LimitReader.Read function will return io.EOF.

Thanks for your time!

Happy coding.