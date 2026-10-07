# [Tutorial: Getting started with generics](https://go.dev/doc/tutorial/generics)


- Generics allows you to define functions or types that can be used with multiple data types.

- With generics, you can write reusable, type-safe functions or types without duplicating their logic or relying on `any` (runtime type assertion). 

- Generics were introduced in Go 1.18 and significantly updated in Go 1.27


## Create a project folder for this tutorial: 

```shell 
mkdir hello-generics
cd hello-generics
```

## Issue with non-generic functions

Suppose you want to define a function that returns the sum of all values in a map. Now based on the type of map, you would need to define multiple functions:

1. For a map of `int64` values:
    ```go
    func SumInts(m map[string]int64) int64 {
        var sum int64
        for _, v := range m {
            sum += v
        }
        return sum
    }
    ```

2. For a map of `float64` values:
    ```go
    func SumFloats(m map[string]float64) float64 {
        var sum float64
        for _, v := range m {
            sum += v
        }
        return sum
    }
    ```

Note, the implementation is mostly same in both functions. With generics, you can club these functions into a single generic function.

## Generic function

In this section, we will define a single generic function that can operate on both `int64` and `float64` value maps.

To make a function generic, we define its **type parameters**.
Basic syntax of a generic function:

```go
func FunctionName[T Constraint](param T) T {}
```

Here `T` is the type parameter name and `Constraint` is the **type constraint** defining what types are allowed and what operations (e.g `+`, `==`, method calls) can be performed on it.

Create a `main.go` file and do following in it:

- Define the generic function:
    - Signature:
        ```go
        func Sum[K comparable, V int64 | float64](m map[K]V) V {}
        ```

        - `comparable` is a built-in type constraint. It represents the set of all types that can be compared using the equality operators `==` and `!=` such as `bool`, numeric types like `int`, pointers, etc.

        We are using `comparable` here since Go requires `map` keys to be comparable. By constraining `K` type param to `comparable`, we can use `K` as the key type for the `m` parameter.

        - `V` is constrained as union of `int64` and `float64`. This will allow the caller to pass either `int64` or `float64` value maps.

    - Body: 
        ```go
        var sum V
        for _, v := range m {
            sum += v
        }
        return sum
        ```

- Call the generic function in main function:
    ```go
    func main() {
        mi := map[string]int64{"a": 1, "b": 2}
        mf := map[string]float64{"a": 1.1, "b": 2.2}
        fmt.Println(Sum[string, int64](mi)) // same as fmt.Println(Sum(mi))
        fmt.Println(Sum[string, float64](mf)) // same as fmt.Println(Sum(mf))
    }
    ```
    
    In this code, type arguments are not needed. Go compiler can infer them from the map argument. 

    But in some cases, the type arguments are necessary. For example: A function without any parameters.


- To simplify the code, we can declare the type constraints separately:
    - Declare a custom type constraint `Number` outside the generic function:
        ```go
        type Number interface {
            int64 | float64
        }
        ```

    - Add a new generic function using the predefined type constraint:
        ```go
        func SumNum[K comparable, V Number](m map[K]V) V {
            return Sum(m)
        }
        ```

- Call the generics in main function:
    ```go
    fmt.Println(SumNum(mi))
    fmt.Println(SumNum(mf))
    ```

- Run the code:

    ```shell
    go run main.go
    ```

