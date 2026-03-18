package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	c "github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/controllers"
)

func main() {
	text := "MoviesStreamReview"
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello, %s", text)
	})

	router.GET("/movies", c.GetMovies())
	router.GET("/movies/:imdb_id", c.GetMovie())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("failed to start server", err)
	}
}