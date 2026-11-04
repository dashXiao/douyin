package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"tiktok/pkg/constants"
)

type Claims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

func CreateToken(userId int64) (string, error) {
	expireTime := time.Now().Add(24 * 7 * time.Hour) // Tokens expire after seven days.
	nowTime := time.Now()                            // Current time
	claims := Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // Expiration timestamp
			IssuedAt:  nowTime.Unix(),    // Issued-at timestamp
			Issuer:    "tiktok",          // Token issuer
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenStruct.SignedString([]byte(constants.JWTValue))
}

func CheckToken(token string) (*Claims, error) {
	response, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JWTValue), nil
	})

	if err != nil {
		return nil, err
	}

	if resp, ok := response.Claims.(*Claims); ok && response.Valid {
		return resp, nil
	}

	return nil, err
}
