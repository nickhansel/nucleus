package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GenerateAccessToken(user_id int32) (string, error) {

	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["type"] = "access"
	claims["user_id"] = user_id
	// exp date that expires in 1 hour
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateRefreshToken(user_id int32) (string, error) {

	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["type"] = "access"
	claims["user_id"] = user_id
	// exp date that expires in 30 days
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func TokenValid(c *gin.Context) error {

	errs := godotenv.Load("../.env")

	if errs != nil {
		fmt.Println("Error loading .env file")
	}

	// validate the token by getting the token string from the header
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

func ExtractToken(c *gin.Context) string {
	// check to make sure the string contains spaces

	// get the bearer token from the Authorization header
	bearerToken := c.GetHeader("Authorization")

	if bearerToken == "" {
		return ""
	}

	// split the string into an array and if there is an error return an empty string
	strArr := strings.Split(bearerToken, " ")

	if (strArr[0] != "Bearer") || (len(strArr) != 2) {
		return ""
	}
	// check if the token is valid
	return strArr[1]
}

func ExtractTokenID(c *gin.Context) (int32, error) {

	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return int32(uid), nil
	}
	return 0, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return password == hash
}
