# Markdown

## Explaining

In this reading, we want to see how io.TeeReader and io.Copy works so i prepared some example for you.

We will download a gz file and we ll save it while doing it we will display how much the file has been downloaded so lets begin.

```go
res, err := http.Get("http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-all-5gram-20120701-0.gz")
if err != nil {
	panic(err)
}
```

First of all, we will download file.

```go
local, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY, 0600)
if err != nil {
	panic(err)
}
defer local.Close()
```

after that we will open a file to write in it.

```go
dec, err := gzip.NewReader(res.Body)
if err != nil {
	panic(err)
}
```

Now we open the gzip file

```go
if _, err := io.Copy(local, io.TeeReader(dec, &progress{})); err != nil {
	panic(err)
}
```

and finally we will copy content of gzip to file we opened but there is something i didnt tell you which is progress of downloaded file. We are gonna handle it with io.TreeReader before dive in to counter{} structure lets examine io.TeeReader func.

```go
// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func TeeReader(r Reader, w Writer) Reader {
	return &teeReader{r, w}
}

type teeReader struct {
	r Reader
	w Writer
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}
```

TeeReader function actually as explained 

> returns a Reader that writes to w what it reads from r
> 

so how are we gonna use it?

While we reading the gzip we will print progress to the screen with the writer i sent to the TeeReader function , so we should write a function that implements a writer interface.

 

```go
type progress struct {
	total uint64
}

func (p *progress) Write(b []byte) (int, error) {
	p.total += uint64(len(b))
	progress := float64(p.total) / (constants.MB)

	fmt.Printf("Downloading %f MB... \n", progress)
	return len(b), nil
}
```

As you can see we created a counter struct and write func which implements io.Writer interface to send it to TeeReader. If you look at the code I wrote 2 ago you can see we are sending counter struct to TeeReader func. 

## Debugging

Before starting, i want to say that I put debug points on the functions below.

- `(f *File) Write(b []byte) (n int, err error)`  (os package)
- `(z *Reader) Read(p []byte) (n int, err error)`  (compress/gzip package)
- `copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error)` (io package)
- `(t *teeReader) Read(p []byte) (n int, err error)` (io package)
- `(c *counter) Write(b []byte) (int, error)` (our main package)

lets start.

When i run the program firsth breakpoint that hit is `copyBuffer` function.

copyBuffer function has been called in io.Copy function as you see.

```go
func Copy(dst Writer, src Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}
```

```go
// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
```

In copyBuffer function we run till first line of for loop with 32KB size buffer and we called src.Read func. 

### src.Read

If we run again we will encounter with

```go
func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}
```

As you can see, here we are calling read func and write func that we gave teereader before.

Lets run program again.

Now we hit gzip read function. If you remember we gave it when we run TeeReader function as `io.Reader`

```go
// Read implements io.Reader, reading uncompressed bytes from its underlying Reader.
func (z *Reader) Read(p []byte) (n int, err error) {
	if z.err != nil {
		return 0, z.err
	}

	n, z.err = z.decompressor.Read(p)
	z.digest = crc32.Update(z.digest, crc32.IEEETable, p[:n])
	z.size += uint32(n)
	if z.err != io.EOF {
		// In the normal case we return here.
		return n, z.err
	}

	...
}
```

After decompress data in gzip it will return us number of readed bytes now lets run again.

and we gave that function to TeeReader function as io.Writer.

```go
func (p *progress) Write(b []byte) (int, error) {
	p.total += uint64(len(b))
	progress := float64(p.total) / (constants.MB)

	fmt.Printf("Downloading %f MB... \n", progress)
	return len(b), nil
}
```

With our write function we are able to print the progress each we read gzip file.

### dst.Write

Now we will call Write function in copyBuffer function which is Write function in os package.

```go
// Write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (f *File) Write(b []byte) (n int, err error) {
	if err := f.checkValid("write"); err != nil {
		return 0, err
	}
	n, e := f.write(b)
	if n < 0 {
		n = 0
	}
	if n != len(b) {
		err = io.ErrShortWrite
	}

	epipecheck(f, e)

	if e != nil {
		err = f.wrapErr("write", e)
	}

	return n, err
}
```

Finally we byte we gave to the file we specified which is file.txt.

and last output is

 

```go
...
...
Downloading 328.406250 MB... 
Downloading 328.437500 MB... 
Downloading 328.468750 MB... 
Downloading 328.500000 MB... 
Downloading 328.531250 MB... 
Downloading 328.562500 MB... 
Downloading 328.593750 MB... 
Downloading 328.625000 MB... 
Downloading 328.656250 MB... 
Downloading 328.687500 MB... 
Downloading 328.718750 MB... 
Downloading 328.750000 MB... 
Downloading 328.781250 MB... 
Downloading 328.812500 MB... 
Downloading 328.843750 MB... 
Downloading 328.869534 MB...
```

Each we read gzip file we print the progress.

Thanks for your time.

Happy coding!