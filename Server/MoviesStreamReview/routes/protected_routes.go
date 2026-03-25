package routes

import (
	"github.com/gin-gonic/gin"
	c "github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/controllers"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.Use(middleware.AuthMiddleware())

	router.GET("/movies/:imdb_id", c.GetMovie(client))
	router.POST("/movies", c.AddMovie(client))
	router.GET("/movies/recommended", c.GetRecommendedMovies(client))
	router.PATCH("/movies/updatereview/:imdb_id", c.AdminReviewUpdate(client))
}
