# [Accessing a relational database](https://go.dev/doc/tutorial/database-access)

The [`database/sql`](https://pkg.go.dev/database/sql) package of the standard library includes types and functions for working with databases.

In this tutorial, we will be using Postgres as our RDBMS.
For detailed documentation on databases, refer [Accessing relational databases](https://go.dev/doc/database/).

## Creating a project

- Create a project folder `hello-db`
    ```shell
    mkdir hello-db
    cd hello-db
    ```
- Create a module for the project:
    ```shell
    go mod init hello_db
    ```

## Setting up a database

- Login in to the database using Postgres CLI `psql`:

    ```shell
    psql -U postgres
    ```

    - `-U` option is used for logging in by username specified in the argument. Without it, the CLI uses the Windows User account username.
    - `postgres` is the default administrator user created during the Postgres installation.
    - On running on the command, you will be prompted for the password of the `postgres` user.
    - On entering the correct password, you will be connected to the `postgres` database if it exists.  
      By default, when the `-d` option is not used, the CLI tries to connect to the database with same name as the username specified in the arguments.

        ```
        PS C:\Users\user> psql -U postgres
        Password for user postgres:

        psql (18.3)
        Type "help" for help.

        postgres=#
        ```

        The phrase `postgres=#`, indicates us about 3 different things:
        1. First part `postgres` means you are connected to the `postgres` database.
        2. `=` means the CLI is ready for a new command. Sometimes you will see `-` instead of `=`,
           which means CLI is still collecting input for the current SQL statement.
        3. `#` indicates that you're connected as a PostgreSQL superuser (or a role with the relevant superuser status). For a normal PostgreSQL user, you will see `>` instead of `#`.

    - To connect with a specific database during login, use `-d` option:
        ```shell
        psql -U postgres -d database_name
        ```

- Create a new database:

    ```sql
    create database golearn;
    ```

    - `create database golearn;` is a SQL statement so it must end with a semi-colon.

- Connect to the `golearn` database:

    ```
    \connect golearn
    ```

    - `\connect` or its short form `\c` is a `psql` meta-command, not SQL. The CLI recognizes it well so it doesn't require a terminating semi-colon.

- Create a SQL script `album.sql` inside the `hello-db` folder:

    ```sql
    DROP TABLE IF EXISTS album;

    CREATE TABLE album (
        id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        title      VARCHAR(128) NOT NULL,
        artist     VARCHAR(255) NOT NULL,
        price      DECIMAL(5,2) NOT NULL
    );

    INSERT INTO album
        (title, artist, price)
    VALUES
        ('Blue Train', 'John Coltrane', 56.99),
        ('Giant Steps', 'John Coltrane', 63.99),
        ('Jeru', 'Gerry Mulligan', 17.99),
        ('Sarah Vaughan', 'Sarah Vaughan', 34.98);
    ```

    This scripts creates the `album` table and seeds it with default data.

- Run the SQL script inside the `psql` CLI:

    ```
    \include ./album.sql
    ```

    - `\include` or `\i` reads and executes the given SQL script.

- Run the `SELECT` statement to verify:
    ```
    golearn=# select * from album;
    1 | Blue Train    | John Coltrane  | 56.99
    2 | Giant Steps   | John Coltrane  | 63.99
    3 | Jeru          | Gerry Mulligan | 17.99
    4 | Sarah Vaughan | Sarah Vaughan  | 34.98
    ```

## Finding and importing a database driver

The Go standard library's [`database/sql`](https://pkg.go.dev/database/sql) package provides a generic, database-agnostic interface for working with databases (connections, transactions, and queries).

However, `database/sql` does not talk directly to any specific database engine over the network.

It relies on an external **database driver** to translate generic SQL calls into the protocol the database understands. A list of community-supported drivers is maintained on the [Go SQLDrivers Wiki](https://go.dev/wiki/SQLDrivers).

### Driver for Postgres

For PostgreSQL, **[pgx](https://pkg.go.dev/github.com/jackc/pgx/v5)** is the modern and recommended driver.

- By default, `pgx` is not just a driver—it is a full-featured,
  standalone PostgreSQL client with its own connection types (`pgx.Conn`, `pgxpool.Pool`).  
  Developers can use native `pgx`, without needing Go's `database/sql package` at all.

- `pgx` includes the [stdlib adapter](https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib),
  which lets you use `pgx` seamlessly with Go's standard `database/sql` package.  
  The package name for the adapter is `github.com/jackc/pgx/v5/stdlib`.

In this tutorial, we will use the `stdlib` adapter:

- Create `main.go` file in `hello-db` folder with following contents:

    ```go
    package main

    import (
        "database/sql"

        "github.com/jackc/pgx/v5"
        "github.com/jackc/pgx/v5/stdlib"
    )
    ```

## Connecting to the database

- Update `main.go` with following code:
    - Declare a package-level variable `db` for storing the [database handle](https://go.dev/doc/database/open-handle).

        ```go
        var db *sql.DB
        ```

        In production, avoid declaring it globally. [Learn more](https://www.alexedwards.net/blog/organising-database-access).

    - Create the `main` function and initialize the connection URI string in it:
        ```go
        connStr := "postgres://username:password@localhost:5432/golearn?sslmode=disable"
        ```
        `sslmode=disable` is used here for local development. Don't use it in production.
    - Call `sql.Open` with driver name `pgx` to get the database handle and assign it to `db` variable:
        ```go
        db, err = sql.Open("pgx", connStr)
        ```
    - Call `db.ping` to confirm that the connection to the PostgreSQL server can actually be established.
        ```go
        err = db.Ping()
        fmt.Print("Connected!")
        ```

## Selecting single row

- Declare an `Album` struct to hold a row of the table:

    ```go
    type Album struct {
        ID     int64
        Title  string
        Artist string
        Price  float32
    }
    ```

- Add a function to get an album row by Id:
    - Function signature:

        ```go
        func albumById(id int) (Album, error)
        ```

    - Function body:
        - Create an instance of `Album` to hold the `album` table row:

            ```go
            var a Album
            ```

        - Execute the SELECT query:
            ```go
            row := db.queryRow("SELECT * from album where id = $1", id)
            ```
            `db.QueryRow` is designed to execute queries which are expectedto return only single row.  
             It doesn't return any error. The error (if any) is deferred until the returned `row` is scanned using `row.Scan`.
        - Scan the row into the `Album` instance:

            ```go
            err := row.Scan(&a.Id, &a.Title, &a.Artist, &a.Price)
            ```

        - Check for `sql.ErrNoRows` error before returing the result:
            ```go
            if err != nil {
                if err == sql.ErrNoRows {
                    return a, fmt.Errorf("albumsByID %d: no such album", id)
                }
                return a, fmt.Errorf("albumsByID %d: %v", id, err)
            }
            return a, nil
            ```

- Call the `albumByID` in `main` function and print the result:

    ```go
    	id := 1
        a, err := albumByID(id)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("albumByID(%d): %v", id, a)
    ```

- Run the code to check the output

## Selecting multiple rows

- Add a function to get multiple albums by the artist name:
    - Function signature:

        ```go
        func albumsByArtist(artist string) ([]Album, error) {}
        ```

    - Function body:
        - Declare a slice of albums:

            ```go
            var albums []Album
            ```

        - Execute the query:

            ```go
            rows, err := db.Query("SELECT * from album where artist = $1", artist)
            ```

            We use `db.Query` here, since it is designed to execute queries which are expected to return multiple rows.

        - Check for error:
            ```go
            if err != nil {
                return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
            }
            ```
        - Use `defer rows.Close()` to ensure connection resources held by the result set are freed when the function returns.

        - Loop through the rows and add each scanned row to the `albums` slice:

            ```go
            for rows.Next(){
                var a Album

                // copy the column values of current row into the struct fields
                err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)

                if err != nil {
                    return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
                }

                albums = append(albums, a)
            }
            ```

        - Call `rows.Err()` after the loop to check whether any error occurred during iteration:

            ```go
            err2 := rows.Err()
            if err2 != nil {
                return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
            }

            // finally return the albums
            return albums, nil
            ```

- Update `main` to call `albumsByArtist`:
    ```go
    artist := "John Coltrane"
    albums, err := albumsByArtist(artist)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("albumsByArtist(%v): %v\n", artist, albums)
    ```

## Inserting new rows

- Add a function to add albums:
    - Function signature:

        ```go
        func addAlbum(album Album) (int, error) {}
        ```

    - Function body:
        - Execute the INSERT query:

            ```go
            row := db.QueryRow(
                "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
                album.Title, album.Artist, album.Price
            )
            ```

            We are using `db.QueryRow` since the query contains the `RETURNING` clause with returns the inserted row.  
             Using `RETURNING id` clause we get the inserted row ID. This clause is supported in PostgreSQL.

        - Scan the row using `row.Scan` to get the inserted row ID:

            ```go
                var id int
                err := row.Scan(&id)
            ```

        - Check for errors:

            ```go
            if err != nil {
                return 0, fmt.Errorf("addAlbum: %v", err)
            }

            return id, nil
            ```

- Update `main` to call `addAlbum`:
    ```go
    aId, err := addAlbum(Album{Title: "Hurry Up Tomorrow", Artist: "Weeknd", Price: 69.69})
    if err != nil {
    	log.Fatal(err)
    }
    fmt.Printf("addAlbum: New album inserted having ID = %v\n", aId)
    ```

## Running the code

Run the code from the `hello-db` project directory to check the output:

```shell
go run .
```
