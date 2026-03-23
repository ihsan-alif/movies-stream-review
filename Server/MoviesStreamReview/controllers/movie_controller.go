package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/database"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/models"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/utils"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/googleai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")
var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch movies"})
		}

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode movies"})
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movie ID is required"})
			return
		}

		var movie models.Movie

		if err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movie models.Movie

		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to add movie"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieId = c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movie id required"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}

		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting review ranking"})
			return
		}

		filter := bson.M{"imdb_id": movieId}
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}

		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to update movie"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}

		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)
	}
}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID not found in context"})
			return
		}

		favouriteGenres, err := GetUserFavouriteGenres(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = godotenv.Load(".env")
		if err != nil {
			log.Println("warning: .env file not found")
		}

		var recommendedMovieLimitVal int64 = 5
		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")

		if recommendedMovieLimitStr != "" {
			recommendedMovieLimitVal, err = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		findOpts := options.Find().SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}}).SetLimit(recommendedMovieLimitVal)

		filter := bson.M{
			"genre.genre_name": bson.M{
				"$in": favouriteGenres,
			}}

		cursor, err := movieCollection.Find(ctx, filter, findOpts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to fetch recommended movies"})
			return
		}

		var recommendedMovies []models.Movie
		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}

		c.JSON(http.StatusOK, recommendedMovies)

	}
}

func GetUserFavouriteGenres(userId string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}

	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)

	var result bson.M
	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return []string{}, err
	}

	favGenresArray, ok := result["favourite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("unable to retrieve favourite genres for user")
	}

	var genreNames []string
	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil
}

func GetReviewRanking(adminReview string) (string, int, error) {
	rankings, err := GetRankings()
	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited += ranking.RankingName + ","
		}
	}

	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	err = godotenv.Load(".env")
	if err != nil {
		log.Println("warning: unable to find .env file")
	}

	googleAi := os.Getenv("GOOGLEAI_API_KEY")
	if googleAi == "" {
		return "", 0, errors.New("could not read GOOGLEAI_API_KEY")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	llm, err := googleai.New(ctx, googleai.WithAPIKey(googleAi), googleai.WithDefaultModel("gimini-1.5-flash"), googleai.WithDefaultTemperature(0.1), googleai.WithDefaultMaxTokens(10))
	if err != nil {
		return "", 0, err
	}

	basePromptTemplate := os.Getenv("BASE_PROMPT_TEMPLATE")
	if basePromptTemplate == "" {
		return "", 0, errors.New("could not read BASE_PROMPT_TEMPLATE")
	}

	basePrompt := strings.Replace(basePromptTemplate, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(ctx, basePrompt+adminReview)
	if err != nil {
		return "", 0, err
	}

	rankVal := 0
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil

}

func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil

}
