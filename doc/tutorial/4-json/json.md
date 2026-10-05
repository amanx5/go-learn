# [JSON](https://go.dev/doc/tutorial/json)

Below are some of the main Go packages for working with JSON:

1. [encoding/json/v2](https://pkg.go.dev/encoding/json/v2) - Newer package with a cleaner API and better defaults
2. [encoding/json](https://pkg.go.dev/encoding/json) - Standard Go v1 package
3. [encoding/json/jsontext](https://pkg.go.dev/encoding/json/jsontext)

## Encoding

Encoding is the process of serializing in-memory Go data structures into JSON text.

To encode JSON data, we use the [Marshal](https://pkg.go.dev/encoding/json/v2#Marshal) function.

### Signature of Marshal function:

```go
func Marshal(in any, opts ...Options) (out []byte, err error)
```

Basic semantics around the `in` argument:

- It must be a valid, non-cyclic data structure that can be represented as JSON.
- If it is a `nil` pointer, it will be encoded as `null`.
- If it is a non-`nil` pointer, the result will be the encoding of the value it points to.
- If `in` is a Go map, the key type must encode as a string: either a `string`, an integer type (which Go converts to a string), or a type implementing `encoding.TextMarshaler` / `encoding.TextAppender`.
- It cannot be a channel, function, or complex number (`complex64`/`complex128`).
- If it is a `struct`, only exported fields (names starting with a capital letter) will be encoded.

### Example of encoding:

- Import the json v2 package

    ```go
    import "encoding/json/v2"
    ```

- For this example, we will be encoding a struct type, so lets declare a struct type named `Message`:

    ```go
    type Message struct {
        Name string
        Body string
        Time time.Time // import time package for time type
    }
    ```

- Create a `Message` instance

    ```go
    msg := Message{"Aman", "Hello", time.Date(2026, 9, 28, 4, 0, 0, 0, time.UTC)}
    ```

- Encode `msg` to JSON bytes using `json.Marshal`:

    ```go
    encoded, err := json.Marshal(msg)
    ```

- If the encoding was successful, then `fmt.Println` would print as commented:
    ```go
    fmt.Println(err)	// <nil>
    fmt.Println(string(encoded))	// {"Name":"Aman","Body":"Hello","Time":"2026-09-28T04:00:00Z"}
    ```

## Decoding

To decode JSON data we use the [Unmarshal](https://go.dev/pkg/encoding/json/v2#Unmarshal) function.

### Signature of Unmarshal function:

```go
func Unmarshal(in []byte, out any, opts ...Options) (err error)
```

### Example of decoding (continuing to [previous example](#example-of-encoding)):

- Declare a placeholder for the decoded data
    ```go
    var msg2 Message
    ```
- Call `Unmarshal` with the `encoded` bytes and the address of `msg`:

    ```go
    err2 := json.Unmarshal(encoded, &msg2)
    ```

- If `encoded` is a valid JSON byte array that can be written into address of `msg`, then after calling Unmarshal, then `fmt.Println` would print as commented:
    ```go
    fmt.Println(err2)	// <nil>
    fmt.Println(msg2)	// {Aman Hello 2026-09-28 04:00:00 +0000 UTC}
    ```

### Decoding arbitrary data

If we don't know the structure of the JSON in advance, and still want to decode it, we can decode it in `any` type variable:

- Consider JSON:
    ```go
    const book = `{
        "Id": 1,
        "Name": "Cosmos",
        "Keywords": ["science", "fiction"]
    }`
    ```
- Unmarshal into `any` type variable:

    ```go
    var result any
    json.Unmarshal([]byte(book), &result)
    ```

- Check the decoded result type and value:
    ```go
    fmt.Printf("\n%T", result) // map[string]interface {}
    fmt.Printf("\n%v", result) // map[Id:1 Keywords:[science fiction] Name:Cosmos]
    ```

Notice, the type of decoded value is a map of keys as string and values as `any` (`interface{}`).

- To work on individual fields, we can iterate the map keys but the `result` type is still `any`, so we must assert its type as map first:

    ```go
    resultMap := result.(map[string]any)
    ```

- Now we can iterate the map keys:

    ```go
    for key,value := range resultMap {

        switch val := value.(type) {
           case string:
           fmt.Println(key, "is string", val)
        }
    }
    ```

## Reference Types

This section will give you an idea about how Unmarshal deals with the fields having reference types such as pointers, slices and maps.

### Slice

In the [previous example](#decoding-arbitrary-data), we learnt how Unmarshal decodes the JSON into a output variable of type `any`.

Now lets decode the same JSON into a variable of type struct:

- Define `Book` struct type:

    ```go
    type Book struct {
        Id int
        Name string
        Keywords []string
    }
    ```

- Create a `Book` instance:

    ```go
    var b Book
    ```

    Since no value was passed in initialization, `b.Keywords` would be a nil slice.

    `b`'s header would contain a null data pointer (0x0), meaning no backing array exists yet.

- Decode the `book` JSON

    ```go
    json.Unmarshal([]byte(book), &b)
    ```

    `Unmarshal` detects that `b.Keywords` is nil slice.  
     To populate the field, `Unmarshal` allocates a new backing array on the heap behind the scenes and updates `b.Keywords`'s slice header (Data pointer, len, and cap).

    In case the `book` json did not had some value for `Keywords` field, then the `b.Keywords` would be left as a nil slice.

### Pointer

Pointer fields are handled in similar way as slices by the Unmarshal function.

Suppose the JSON to decode is of type `Foo`:

```go
type Foo struct {
    Bar *Bar
}

type Bar struct {
    Size int
}
```

If the JSON text passed to Unmarshal contains the Bar field of type `Bar`, then the Unmarshal would allocate a new `Bar` and populate it in the given instance of `Foo`, otherwise `Foo.Bar` would remain as a nil pointer.

## Streaming Marshal and Unmarshal

The `encoding/json/v2` package (requires Go 1.27) provides `MarshalWrite` and `UnmarshalRead` functions for working with streams directly without requiring an intermediate `[]byte`.


### Marshal vs MarshalWrite

While `json.Marshal` returns an in-memory byte slice,
`json.MarshalWrite` writes the JSON data directly to an `io.Writer` stream without allocating an intermediate byte slice.

**Signatures:**

```go
func Marshal(in any, opts ...Options) (out []byte, err error)
func MarshalWrite(out io.Writer, in any, opts ...options) (err error)
```

`MarshalWrite` combines marshaling and streaming into a single, straightforward operation.

### Example of MarshalWrite function

- Suppose we have an empty file `message.json` in which we want to store JSON data of following structure:
    ```go
    type Message struct {
        Body string
        Time int
    }
    ```
- Open the file (with write perms) in the program:
    ```go
    file, err := os.OpenFile(name, os.O_RDWR, 0666)
    ```
- Create a `Message` instance with desired JSON data:
    ```go
    var msg = Message{Body: "Hello", Time: int(time.Now().Unix())}
    ```
- Write the `msg` directly to the file:
    ```go
    json.MarshalWrite(file, msg)
    ```

### Unmarshal vs UnmarshalRead

The `json.Unmarshal` processes data already loaded in-memory as a byte slice (`[]byte`),
whereas `json.UnmarshalRead` streams data directly from an I/O source (`io.Reader`)

- Use `json.Unmarshal` when you already have a byte slice in hand (e.g., pulling a payload from a Redis cache, database blob, or local string literal).
- Use `json.UnmarshalRead` whenever data originates from a stream (e.g., web requests, reading files from disk, or network sockets)

**Signatures:**

```go
func Unmarshal(in []byte, out any, opts ...Options) (err error)
func UnmarshalRead(in io.Reader, out any, opts ...options) (err error)
```

### Example of UnmarshalRead function

- After running the [MarshalWrite example](#example-of-marshalwrite-function), the `message.json` file would have content like below:
    ```json
    { "Body": "Hello", "Time": 1791218035 }
    ```
- Open the file in the program:
    ```go
    file, err = os.Open(filePath)
    ```
- Create an `Message` instance for storing the file contents:
    ```go
    var msg Message
    ```
- Read the file contents and marshal it simultaneously into the `msg` variable:
    ```go
    json.UnmarshalRead(file, &msg)
    ```

## Custom Marshalling

Suppose we have a type representing Semantic version numbers:
```go
type Semver struct {
	Major, Minor, Patch int
}
```

By default, marshalling `Semver{0, 0, 1}` would output `{"Major":0,"Minor":0,"Patch":1}`
```go
out1, _ := json.Marshal(Semver{0, 0, 1})
fmt.Println(string(out1))    // {"Major":0,"Minor":0,"Patch":1}
```

To customize marshalling behaviour of `Semver` type:

-   The `Semver` type must implement the `json.MarshalerTo` interface:
    ```go
    type MarshalerTo interface {
        MarshalJSONTo(*jsontext.Encoder) error
    }
    ```
-   To implement `json.MarshalerTo` interface, the `Semver` type must have MarshalJSONTo method:
    ```go
    func (sv Semver) MarshalJSONTo(enc *jsontext.Encoder) error {
        return json.MarshalEncode(enc, fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch))
    }
    ```
-   Marshal will run the custom marshalling logic if the `Semver` type correctly implements the 
    `json.MarshalerTo` interface and would output `"0.0.1"`:
    ```go
	out2, _ := json.Marshal(Semver{0, 0, 1})
	fmt.Println(string(out2)) // "0.0.1"
    ```

## Custom Unmarshalling

Suppose we want to decode a semantic version string `"0.0.1"` into our `Semver` struct.

By default, unmarshalling `"0.0.1"` into `Semver` would fail because Go expects a JSON object (`{...}`), not a JSON string:

```go
var sv1 Semver
err1 := json.Unmarshal([]byte(`"0.0.1"`), &sv1)
fmt.Println(err1) // json: cannot unmarshal JSON string into Go main.Semver
```

To customize unmarshalling behaviour of Semver type:

-   The `Semver` type must implement the `json.UnmarshalerFrom` interface:
    ```go
    type UnmarshalerFrom interface {
        UnmarshalJSONFrom(*jsontext.Decoder) error
    }
    ```
-   To implement `json.UnmarshalerFrom` interface, the pointer type `*Semver`
    must have an UnmarshalJSONFrom method (a pointer receiver is required
    so it can mutate the struct's fields):
    ```go
    func (sv *Semver) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
        var str string
        err := json.UnmarshalDecode(dec, &str);
        if err != nil {
            return err
        }
        _, err2 := fmt.Sscanf(str, "%d.%d.%d", &sv.Major, &sv.Minor, &sv.Patch)
        return err2
    }
    ```
-   Unmarshal will run the custom unmarshalling logic if the `Semver` type
    correctly implements the `json.UnmarshalerFrom` interface and will populate
    the struct fields:
    ```go
    var sv2 Semver
    err2 := json.Unmarshal([]byte(`"0.0.1"`), &sv2)
    fmt.Println(err2) // <nil>
    fmt.Println(sv2)  // {0 0 1}
    ```

## Quick Comparison of JSON methods

In `encoding/json/v2`, there are 3 pairs of functions targeting different mediums:

| Operation | Target / Source Medium | Functions |
| :--- | :--- | :--- |
| Buffer | In-memory `[]byte` | `json.Marshal` / `json.Unmarshal` |
| Stream | `io.Writer` / `io.Reader` | `json.MarshalWrite` / `json.UnmarshalRead` |
| Token Stream | `*jsontext.Encoder` / `*jsontext.Decoder` | `json.MarshalEncode` / `json.UnmarshalDecode` |


