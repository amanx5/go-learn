package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var db *sql.DB

func main() {
	var err error

	// Get a database handle.
	// [sql.Open] only validates the arguments and may not connect immediately
	db, err = sql.Open("pgx", getConnString())
	if err != nil {
		log.Fatal(err)
	}

	// establish/verify the connection
	// [db.Ping] forces an initial connection attempt.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	id := 1
	alb, err := albumByID(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("albumByID(%d): %v", id, alb)

	artist := "John Coltrane"
	albums, err := albumsByArtist(artist)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("albumsByArtist(%v): %v\n", artist, albums)

	aId, err := addAlbum(Album{Title: "Hurry Up Tomorrow", Artist: "Weeknd", Price: 69.69})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("addAlbum: New album inserted having ID = %v\n", aId)
}

// albumByID queries for the album with the specified ID.
func albumByID(id int) (Album, error) {
	fmt.Println()

	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)

	if err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumByID %d: no such album", id)
		}
		return alb, fmt.Errorf("albumByID %d: %v", id, err)
	}

	return alb, nil
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(artist string) ([]Album, error) {
	fmt.Println()
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * from album where artist = $1", artist)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var a Album

		// copy the column values of current row into the struct fields
		err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)

		if err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
		}

		albums = append(albums, a)
	}

	err2 := rows.Err()
	if err2 != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
	}

	return albums, nil
}

// addAlbum inserts the given album
func addAlbum(album Album) (int, error) {
	fmt.Println()

	row := db.QueryRow(
		"INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		album.Title, album.Artist, album.Price,
	)

	var id int
	err := row.Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}

// this function uses LastInsertId function which is not supported in PostgreSql
func addAlbumMySql(album Album) (int, error) {
	fmt.Println()
	res, err := db.Exec("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", album.Title, album.Artist, album.Price)

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	id, err2 := res.LastInsertId()
	if err2 != nil {
		return 0, fmt.Errorf("addAlbum: %v", err2)
	}

	return int(id), nil
}

func getConnString() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DBUSER")
	pass := os.Getenv("DBPASS")

	// one liner (unreliable if password contains special characters like @, :, /, ?, #, etc. as have special meaning in a connection URL):
	// connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/golearn?sslmode=disable",user, pass)

	// using connection config (escapes special characters and type safe):
	cfg, err := pgx.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}
	cfg.User = user
	cfg.Password = pass
	cfg.Host = "localhost"
	cfg.Port = 5432
	cfg.Database = "golearn"

	connStr := stdlib.RegisterConnConfig(cfg)
	return connStr
}
