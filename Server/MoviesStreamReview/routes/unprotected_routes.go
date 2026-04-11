package routes

import (
	"github.com/gin-gonic/gin"
	c "github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUnProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.GET("/movies", c.GetMovies(client))
	router.GET("/movies/genres", c.GetGenres(client))
	router.POST("/users/register", c.RegisterUser(client))
	router.POST("/users/login", c.LoginUser(client))
	router.POST("/users/logout", c.LogoutUser(client))
	router.POST("/users/refresh", c.RefreshTokenHandler(client))
}