# Streaming

The term streaming in this document refers to process of consuming a sequence of bytes coming from a source incrementally.

## Reader interface

The [`io.Reader`](https://pkg.go.dev/io#Reader) interface represents a **source** from which bytes can be consumed using its `Read` method.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

The source can be:

- an opened file
    ```go
    source := os.Open() // *os.File
    ```
- the process's standard input
    ```go
    source := os.Stdin  // *os.File
    ```
- an HTTP response body
    ```go
    source := resp.Body // *io.ReadCloser
    ```
- a reader over a string in memory
    ```go
    source := strings.NewReader("Aman") // *strings.Reader
    ```
- a reader over a byte slice in memory
    ```go
    source := bytes.NewReader(byte[]("Aman")) // *bytes.Reader
    ```
- an in-memory growable byte buffer
    ```go
    source := &bytes.Buffer{} // *bytes.Buffer
    ```

All of these sources have a `Read` method, hence they satisfies the `Reader` inferface.

### Reading an HTTP response body:

- Initiate a HTTP GET request:
    ```go
    resp, _ := http.Get("https://dummyjson.com/users/?limit=1")
    ```
- Create a buffer that will be used for storing the reponse body:
    ```go
    buf := make([]byte, 4096)
    ```
- Now, the web server might send the data in chunks over time:
    ```
    time ------------------------------------------------→
            ──data──      ──data──       ──data──
                ↓             ↓              ↓
              Read()        Read()         Read()
    ```
- We read the data continously until we find out that there is no more data:
    ```go
    for {
        n, err := resp.Body.Read(buf)

        // process buf[:n]

        if err == io.EOF {
            break
        }
    }
    ```

## Writer interface

The [`io.Writer`](https://pkg.go.dev/io#Writer) interface represents a **destination** in which we can write bytes.

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

Any type that provides a `Write` method with signature below, implements `io.Writer`.

```go
// p []byte refers to data that caller wants to write
// n tells the caller how many bytes were actually written
Write(p []byte) (n int, err error)
```

The destination can be:

- an opened file
    ```go
    dest := os.Open()
    ```
- the process's standard output
    ```go
    dest := os.Stdout
    ```
- a buffer
- a network connection
- a custom API:
    - Consider a API type `Terminal` 

        ```go
        type Terminal struct {}

        // Write method writes the bytes to the terminal
        func (t Terminal) Write(p []byte) (int, error) {
            fmt.Println(string(p))
            return len(p), nil
        }
        ```

    - A `Terminal` instance can be assigned to a `io.Writer` variable:
        ```go
        var w io.Writer = Terminal{}
        w.Write([]byte("hello"))
        ```
