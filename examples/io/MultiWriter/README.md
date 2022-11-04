# MultiWriter

## Explaining

In this reading, we want to see how io.MultiWriter works so i prepared an example for you.

We will create a structure and encode it as xml format after that we will write it to different locations like file or os.stdOut so lets begin.

```go
type Person struct {
	Name    string
	Surname string
}

func (p *Person) Write(b []byte) (n int, err error) {
	fmt.Printf("Write method ran. Person's name is %s \n", p.Name)
	return len(b), nil
}
```

First of all, we will create a structure which implements io.Writer.

```go
persons := []Person{
		{Name: "name1", Surname: "surname1"},
		{Name: "name2", Surname: "surname2"},
		{Name: "name3", Surname: "surname3"},
		{Name: "name4", Surname: "surname4"},
	}
```

after that we create instances of person struct.

```go
	f, err := os.Create("1.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := bytes.NewBuffer(make([]byte, 0))
```

Now we define things which implements io.Writer interface.

```go
	w := io.MultiWriter(f, buf, os.Stdout)
	for i := range persons {
		w = io.MultiWriter(w, &persons[i])
	}
```

We want to write xml data through all io.Writer interfaces when we call Write function. 

```go
	enc := xml.NewEncoder(w)
	if err := enc.Encode(persons); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Buffer string:", buf.String())
```

and finally we will call NewEncoder and send data in it. When we call new encoder it encode data to xml first after calls Write function to write it. When NewEncoder function calls the write function it actually calls the io.multiWriter’s writer function.

```go
type multiWriter struct {
	writers []Writer
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = ErrShortWrite
			return
		}
	}
	return len(p), nil
}

// MultiWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command.
//
// Each write is written to each listed writer, one at a time.
// If a listed writer returns an error, that overall write operation
// stops and returns the error; it does not continue down the list.
func MultiWriter(writers ...Writer) Writer {
	allWriters := make([]Writer, 0, len(writers))
	for _, w := range writers {
		if mw, ok := w.(*multiWriter); ok {
			allWriters = append(allWriters, mw.writers...)
		} else {
			allWriters = append(allWriters, w)
		}
	}
	return &multiWriter{allWriters}
}
```

As you see above, when we call io.MultiWriter function it actually creates io.Writer slice and appends io.Writer interfaces in it. If one of arguments is multiWriter then it appends all of io.Writer interfaces in multiWriter too.

When you call its Write function, it runs all of Write function of all io.Writer interfaces. 

## Debugging

Before starting, i want to say that I put debug points on the functions below.

- `(b *Buffer) Write(p []byte) (n int, err error)` (bytes package)
- `(f *File) Write(b []byte) (n int, err error)`  (os package)
- `(p *Person) Write(b []byte) (n int, err error)` (main package)
- `(enc *Encoder) Encode(v interface{}) error` (xml package)
- `(t *multiWriter) Write(p []byte) (n int, err error)` (io package)

lets start.

When i run the program first breakpoint that hit is `Encode` function.

`Encode` will call MultiWriter’s Write function with `enc.p.Flush()` function.

```go
func (enc *Encoder) Encode(v interface{}) error {
	err := enc.p.marshalValue(reflect.ValueOf(v), nil, nil)
	if err != nil {
		return err
	}
	return enc.p.Flush()
}
```

Lets dive in to `enc.p.Flush()` function.

```go
// Flush writes any buffered data to the underlying io.Writer.
func (b *Writer) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.wr.Write(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}
```

We call Write function with  b.wr.Write(b.buf[0:b.n]) as you see. Now run program again.

### wr.Write()

```go
func (t *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = ErrShortWrite
			return
		}
	}
	return len(p), nil
}
```

Now we will run all Write functions sequentially in a for loop. 

Run program till breakpoint we put.

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

We hit breakpoint which is in os.File.Write because we first appended os.File as io.Writer. 

```go
io.MultiWriter(f, buf, os.Stdout)
```

```go
// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	m, ok := b.tryGrowByReslice(len(p))
	if !ok {
		m = b.grow(len(p))
	}
	return copy(b.buf[m:], p), nil
}
```

Now we hit Write function at bytes package.

after this function ran we hit `os.File.Write` because os.StdOut is also like a file and when we try to write it we call same function as we did.

After we start to call Person Write functions.

```go
func (p *Person) Write(b []byte) (n int, err error) {
	fmt.Printf("Write method ran. Person's name is %s \n", p.Name)
	return len(b), nil
}
```

We will call this function for 4 times and our program will end.

Our programs output is

```go
<Person><Name>name1</Name><Surname>surname1</Surname></Person><Person><Name>name2</Name><Surname>surname2</Surname></Person><Person><Name>name3</Name><Surname>surname3</Surname></Person><Person><Name>name4</Name><Surname>surname4</Surname></Person>Write method ran. Person's name is name1 
Write method ran. Person's name is name2 
Write method ran. Person's name is name3 
Write method ran. Person's name is name4 
Buffer string: <Person><Name>name1</Name><Surname>surname1</Surname></Person><Person><Name>name2</Name><Surname>surname2</Surname></Person><Person><Name>name3</Name><Surname>surname3</Surname></Person><Person><Name>name4</Name><Surname>surname4</Surname></Person>
```

and we created 1.txt. Content of 1.txt is 

```go
<Person><Name>name1</Name><Surname>surname1</Surname></Person><Person><Name>name2</Name><Surname>surname2</Surname></Person><Person><Name>name3</Name><Surname>surname3</Surname></Person><Person><Name>name4</Name><Surname>surname4</Surname></Person>
```

Thanks for your time.

Happy coding!