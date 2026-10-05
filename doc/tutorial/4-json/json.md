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
