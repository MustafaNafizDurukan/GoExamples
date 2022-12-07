# Difference of SectionReader and LimitReader

The **`io.SectionReader`** and **`io.LimitReader`** types are both types of **`io.Reader`** in the Go standard library **`io`** package. They both provide methods for reading only a portion of an underlying **`io.Reader`** without having to read the entire thing into memory. However, they differ in the way they specify the section of the underlying reader to read from.

The **`io.SectionReader`** type specifies the section of the underlying reader to read from using a starting offset and a number of bytes to read. For example, you could create a new **`io.SectionReader`** that reads from the file **`myfile.txt`** starting at byte 10 and ending at byte 30 like this:

```
Copy code
// Import the io package
import "io"

// Create a new io.Reader that reads from a file
file, err := os.Open("myfile.txt")
if err != nil {
    // Handle the error
}

// Create a new io.SectionReader that reads from the file starting at byte 10
// and ending at byte 20
reader := io.NewSectionReader(file, 10, 20)

```

In contrast, the **`io.LimitReader`** type specifies the section of the underlying reader to read from using a maximum number of bytes to read. For example, you could create a new **`io.LimitReader`** that reads from the file **`myfile.txt`** and stops reading after reading 20 bytes like this:

```
Copy code
// Import the io and ioutil packages
import (
    "io"
    "ioutil"
)

// Create a new io.Reader that reads from a file
file, err := os.Open("myfile.txt")
if err != nil {
    // Handle the error
}

// Create a new io.LimitReader that reads from the file and stops after reading 20 bytes
reader := io.LimitReader(file, 20)

```

In this case, the **`io.LimitReader`** will read from the file starting at the current position and stop after reading 20 bytes, regardless of where that position is in the file.

In summary, the main difference between the **`io.SectionReader`** and **`io.LimitReader`** types is the way they specify the section of the underlying reader to read from. The **`io.SectionReader`** uses a starting offset and a number of bytes to read, while the **`io.LimitReader`** uses a maximum number of bytes to read from the current position.