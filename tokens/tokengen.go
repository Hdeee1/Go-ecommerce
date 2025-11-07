package tokens

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type SignedDetail struct {
	Email		string
	First_Name 	string
	Last_Name	string
	User_ID		string
	jwt.RegisteredClaims
}

var SECRET_KEY []byte

func init() {
	godotenv.Load()
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "bukankan_ini_my_secret_key_ku"
	}
	SECRET_KEY = []byte(secretKey)
}

func TokenGenerator(email, firstName, lastName, userID string) (string, string, error) {
	claims := &SignedDetail{
		Email: email,
		First_Name: firstName,
		Last_Name: lastName,
		User_ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedDetail{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
} 

func ValidateToken(SignedToken string) (*SignedDetail, error) {
	token, err := jwt.ParseWithClaims(
		SignedToken,
		&SignedDetail{},
		func(t *jwt.Token) (interface{}, error) {
			return SECRET_KEY, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetail)
	if !ok {
		return nil, err
	}

	return claims, nil
}