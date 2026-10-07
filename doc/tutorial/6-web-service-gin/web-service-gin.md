# [Developing a RESTful API with Go and Gin](https://go.dev/doc/tutorial/web-service-gin)


This tutorial introduces basic of building REST APIs with Go and [Gin Web Framework](https://gin-gonic.com/en/docs/).

## Design API endpoints

We will implement following endpoints in this tutorial:

- `GET /albums`: Get a list of all albums as JSON
- `GET /albums/:id` - Get an album by its ID
- `POST /albums` - Add a new album from request data

## Create project folder

- Create a folder `hello-gin`:

    ```shell
    mkdir hello-gin
    cd hello-gin
    ```

- Initialize module:

    ```shell
    go mod init hello_gin
    ```

## Create data

- Create `main.go` file with following code:
    - Mark the file as standalone program:
        ```go
        package main
        ```

    - Declare the type of album:
        
        ```go
        type Album struct {
            ID string `json:"id"`
            Title string `json:"title"`
            Artist string `json:"artist"`
            Price float64 `json:"price"`
        }
        ```

    - Declare sample album data:
        ```go
        var albums = []Album {
              {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
            {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
            {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},  
        }
        ```

## Implement GET /albums

- Add following code in `main.go`:

    - Import the necessary packages:
        ```go
        import (
            "net/http"

            "github.com/gin-gonic/gin"
        )
        ```
        Run `go mod tidy` to add the `gin` package as dependency.

    - Create handler for the endpoint:
        - The signature of handler would be:
            ```go
            func getAlbums(c *gin.Context) {}
            ```

            `gin.Context` carries request details, validates and serializes JSON and more. This is different from Go
            's built-in [context](https://pkg.go.dev/context) package.


        - Body of handler:
            ```go
            c.IndentedJSON(http.StatusOK, albums)
            ```

            Use [Context.IndentedJSON](https://pkg.go.dev/github.com/gin-gonic/gin#Context.IndentedJSON) method to return a response with 200 status and indented albums json as body. If you don't want to indent the repsonse JSON, you can use [`Context.JSON`](https://pkg.go.dev/github.com/gin-gonic/gin#Context.JSON) method instead. 

    - Declare `main` function and setup the router:
        ```go
        func main () {
            // create a router
            router := gin.Default()

            // attach the handler for endpoint
            router.GET("/albums", getAlbums)

            // attach the router to an http.Server and start it
            router.Run("localhost:8080")
        }
        ```
        - We used [gin.Default](https://pkg.go.dev/github.com/gin-gonic/gin#Default) to initialize a gin router since it already attaches a logger and recovery middleware with itself.

## Implement GET /albums/:id

- Create handler:
    - Signature:
    ```go
    func getAlbumByID(c *gin.Context) {}
    ```
    - Body:
        - Retrieve the ID of the requested album:
            ```go
            id := c.Param("id")
            ```
            `Context.Param` retrieves the given parameter from the request URL params.

        - Find the requested album by id from `albums` list:
            ```go
            for _, a := range albums {
                if a.ID == id {
                    c.IndentedJSON(http.StatusOK, a)
                    return
                }
            }

        - If no album found, return a simple 404 error:
            ```go
            c.String(http.StatusNotFound, "Album not found")
            ```

- Attach handler to the router
    ```go
    router.GET("/albums/:id", getAlbumByID)
    ```

## Implement POST /albums

- Create handler: 
    - Signature:
        ```go
        func postAlbums(c *gin.Context) {}
        ```
    - Body:
        - Create an instance of `Album`:
            ```go
            var a Album
            ```

        - Bind the received JSON to the instance `a`:
            ```go
            err := c.BindJSON(&a)
            ```
            `Context.BindJSON` binds the request payload with the given pointer.
            In case the binding fails, it calls `Context.AbortWithError` automatically and returns the error.

        - Check if the binding failed:
            ```go
            if err != nil {
                return
            }
            ```
        - Append the provided album and return response with 201 status code:
            ```go
            albums = append(albums, a)
            c.String(http.StatusCreated, "Album added")
            ```

- Attach handler to the router:
    ```go
    router.POST("/albums", postAlbums)
    ```

## Test endpoints

- Run the code:
    ```
    go run .
    ```

- Test the API endpoints:
    ```shell
    curl localhost:8080/albums
    ```

    ```shell
    curl localhost:8080/albums/1
     ``` 

    ```shell
    curl localhost:8080/albums/4
    ``` 

    ```shell
    curl localhost:8080/albums --include --header "Content-Type: application/json" --request "POST" --data '{"id": "4","title": "Hello Gin","artist": "Giiner","price": 0}'
    ``` 

    ```shell
    curl localhost:8080/albums/4
    ``` 