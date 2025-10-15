package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "secret"

func GenerateToken(email string, userId int64) (string, error) {
	// claims are the key-value pairs in the payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey)) // encode header and payload, signs with secret key
}

func VerifyToken(token string) (int64, error) {
    // parse token and verify the signing method is HMAC (e.g. HS256)
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		// can perform checks first
        _, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
        // return secret key so jwt.Parse() can verify the signature
		return []byte(secretKey), nil
	})
    if err != nil {
        return 0, errors.New("could not parse token")
    }

    tokenIsValid := parsedToken.Valid
    if !tokenIsValid {
        return 0, errors.New("invalid token")
    }

    // type assertion to make sure claims are of type jwt.MapClaims
    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok {
        return 0, errors.New("invalid token claims")
    }

    // email := claims["email"].(string)
    userId := int64(claims["userId"].(float64)) // jwt encodes numbers in float64 
    return userId, nil
}
