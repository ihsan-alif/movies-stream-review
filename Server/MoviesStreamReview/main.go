package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/routes"
)

func main() {
	text := "MoviesStreamReview"
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello, %s", text)
	})

	routes.SetupProtectedRoutes(router)
	routes.SetupUnProtectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("failed to start server", err.Error())
	}
}
