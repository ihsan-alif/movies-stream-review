package utils

import (
	"context"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/ihsan-alif/movies-stream-review/Server/MoviesStreamReviewServer/database"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var (
	SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))
	SECRET_REFRESH_KEY = []byte(os.Getenv("SECRET_REFERESH_KEY"))
	userCollection = database.OpenCollection("users")
)

func GenerateAllToken(email, firstName, lastName, role, userId string) (string, string, error) {
	claims := &SignedDetails{
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Role: role,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "movies-stream-review",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", "", nil
	}

	refreshClaims := &SignedDetails{
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Role: role,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "movies-stream-review",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*7*time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(SECRET_REFRESH_KEY)
	if err != nil {
		return "", "", nil
	}

	return signedToken, signedRefreshToken, nil
}

func UpdateAllTokens(userId, token, refreshToken string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updatedAt:= time.Now()

	updateData := bson.M{
		"$set": bson.M{
			"token": token,
			"refresh_token": refreshToken,
			"updated_at": updatedAt,
		},
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	if err != nil {
		return err
	}

	return nil
}