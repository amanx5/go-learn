# [Return greetings for multiple people](https://go.dev/doc/tutorial/greetings-multiple-people)


In this tutorial, we will add a new function that receives names of multiple people and greets each one of them in a different language.


1. Add a new function `GreetMany` in the [greetings](./greetings/greetings.go) package.

    - Write the function signature:

        ```go
        // GreetMany takes a slice of names and returns a map pairing each name
        // with its own randomized greeting message or an error if a name is empty.
        func GreetMany(names []string) (map[string]string, error) {
            // implementation
        }
        ```

    - Write the function body: 
        - Initialize a `messages` map:
            ```go
            // Create a messages map to associate each of the received names (as a key) 
            // with a generated message (as a value). In Go, you initialize a map with 
            // the following syntax: make(map[key-type]value-type).
            messages := make(map[string]string)
            ```

            In above code, `messages` is a [map](https://go.dev/blog/maps) of `string` keys and `string` values. 
            
            Syntax to initialize a map: 
            ```go
            make(map[keyType]valueType)
            ```

        - Loop through the provided `names`:
            ```go
            for _, name := range names {
                // loop logic
            }
            ```

            [`range`](https://go.dev/wiki/Range) along with for loop i a way to iterate slices or maps.  
            When ranging over a slice, two values are returned for each iteration. The first is the index, 
            and the second is a copy of the element at that index.  
            We don't need the index value here, so we used the Go [blank identifier](https://go.dev/doc/effective_go.html#blank)
            (an underscore) to ignore it.


        - Add logic inside the `for` loop:
            Call `Greet` with the current name in iteration. If `Greet` returns an error, return from the loop immediately with the error.
            Otherwise keep saving the messages in the map. 
            ```go
            message, err := Greet(name)
            // Return any error encountered immediately.
            if err != nil {
                return nil, err
            }

            // Save the message in the map.
            messages[name] = message
            ```

        - Outside the loop body, return the `messages` map as first value  and `nil` as second value.
            ```go
            return messages, nil
            ```

2. Edit [hello-people.go](./hello-people/hello.go) 

    - Declare a slice of names:
        ```go
        names := []string {"Aman" , "John", "Adel"}
        ```
    
    - Call `GreetMany` with the names:
        ```go
        messages, err := greetings.GreetMany(names)
        ```

    - Print the `messages` map if `err` was `nil`