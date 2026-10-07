package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	// Initialize a Gin router
	router := gin.Default()

	// Attach handlers
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	//Start the router at `localhost:8080`
	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	a := findAlbumById(id)

	if a != nil {
		c.IndentedJSON(http.StatusOK, a)
		return
	}

	c.String(http.StatusNotFound, "Album not found")
}

func postAlbums(c *gin.Context) {
	var a Album

	// we are using [gin.Context.ShouldBindJSON] here instead of [gin.Context.BindJSON]
	// since ShouldBindJSON allows us to customize the response in case binding fails
	// whereas BindJSON automatically calls [gin.Context.AbortWithError] if binding fails
	err := c.ShouldBindJSON(&a)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid payload: %s", err.Error())
		return
	}

	// check for existing album
	ea := findAlbumById(a.ID)
	if ea != nil {
		c.String(http.StatusBadRequest, "An album with given id already exists")
		return
	}

	// finally add the album
	albums = append(albums, a)
	c.String(http.StatusCreated, "Album added")
}

func findAlbumById(id string) *Album {
	// not using for i,a := range albums // as `a` will be temp loop variable copy, not original albums[i]
	for i := range albums {
		if albums[i].ID == id {
			return &albums[i]
		}
	}

	return nil
}
