package internal

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID int `json:"user_id" binding:"required"`
}

func ParseSessionToken(user_id int, secrete_key string) (string, error) {
	signedtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ud":  user_id,                              // user_id must be stored in token
		"exp": time.Now().Add(2 * time.Hour).Unix(), // expire date
	})

	token, err := signedtoken.SignedString([]byte(secrete_key))
	if err != nil {
		return "", err
	}

	return token, nil
}

func JwtVerifySignature(token_string string, secrete_key string) (*jwt.Token, error) {
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
		return []byte(secrete_key), nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid Token Signature")
	}

	return token, nil
}

func ExtractAuthBearerToken(ctx *fiber.Ctx) (UserClaims, error) {

	authorization := ctx.Get("Authorization")

	bearerAuth := strings.Replace(authorization, "Bearer", "", -1)
	bearerToken := strings.Replace(bearerAuth, " ", "", -1)

	var secretKey string = os.Getenv("SIGN_IN_PRIVATE_KEY")
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return UserClaims{}, errors.New("invalid Token Signature")
		}

		if errors.Is(err, jwt.ErrTokenMalformed) {
			return UserClaims{}, errors.New("invalid Token Signature")
		}
		fmt.Println(err)
		return UserClaims{}, errors.New("ops something went wrong")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("not ok")
		return UserClaims{}, errors.New("token Claims is Not ok")
	}
	var CalimsData UserClaims
	UserID, err := strconv.Atoi(fmt.Sprint(claims["ud"]))
	if err != nil {
		fmt.Println(err)
		return UserClaims{}, err
	}

	CalimsData.UserID = UserID

	return CalimsData, nil
}
