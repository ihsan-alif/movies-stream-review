package routes

import (
	"github.com/gin-gonic/gin"
	c "github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/controllers"
)

func SetupUnProtectedRoutes(router *gin.Engine) {
	router.GET("/movies", c.GetMovies())
	router.POST("/users/register", c.RegisterUser())
	router.POST("/users/login", c.LoginUser())
}