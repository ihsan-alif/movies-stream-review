package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/database"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/models"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check existing user"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to hash password"})
			return
		}

		user.UserId = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, result)

	}
}

func LoginUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var userLogin models.UserLogin

		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		var foundUser models.User

		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		token, refreshToken, err := utils.GenerateAllToken(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		if err := utils.UpdateAllTokens(foundUser.UserId, token, refreshToken, client); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update token"})
			return
		}

		c.JSON(http.StatusOK, models.UserResponse{
			UserId:          foundUser.UserId,
			FirstName:       foundUser.FirstName,
			LastName:        foundUser.LastName,
			Email:           foundUser.Email,
			Role:            foundUser.Role,
			Token:           token,
			RefreshToken:    refreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})
	}
}

func LogoutUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		var UserLogout struct {
			UserId string `json:"user_id"`
		}

		if err := c.ShouldBindJSON(&UserLogout); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
			return 
		}

		if err := utils.UpdateAllTokens(UserLogout.UserId, "", "", client); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error logging out"})
			return 
		}

		c.JSON(http.StatusOK, gin.H{"message": "logout successfully"})

	}
}

func RefreshTokenHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unable to retrieve refresh token from cookie", "details": err.Error()})
			return 
		}

		claim, err := utils.ValidateRefreshToken(refreshToken)
		if err != nil || claim == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token", "details": err.Error()})
			return 
		}

		var userCollection *mongo.Collection = database.OpenCollection("user", client)

		var user models.User

		if err := userCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: claim.UserId}}).Decode(&user); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return 
		}

		newtoken, newRefreshToken, _ := utils.GenerateAllToken(user.Email, user.FirstName, user.LastName, user.Role, claim.UserId)
		err = utils.UpdateAllTokens(claim.UserId, newtoken, newRefreshToken, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating token"})
			return 
		}

		c.JSON(http.StatusOK, gin.H{"message": "tokens refreshed"})

	}
}