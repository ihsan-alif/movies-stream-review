package routes

import (
	"github.com/gin-gonic/gin"
	c "github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/controllers"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/middleware"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	router.GET("/movies/:imdb_id", c.GetMovie())
	router.POST("/movies", c.AddMovie())
}
