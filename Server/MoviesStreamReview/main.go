package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/database"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/routes"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	text := "MoviesStreamReview"
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello, %s", text)
	})

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: unable to find .env file")
	}

	var client *mongo.Client = database.Connect()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("failed to reach server: %v", err)
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("failed to disconnect from mongoDB: %v", err)
		}
	}()

	routes.SetupUnProtectedRoutes(router, client)
	routes.SetupProtectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("failed to start server", err.Error())
	}
}
