package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	Id    int    `json:"user_id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
}

func GenerateJWT(id int, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorMalformed)
	}

	return token, nil
}

func GetUserClaims(token *jwt.Token) (*UserClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.NewValidationError("invalid token claims", jwt.ValidationErrorMalformed)
	}

	userClaims := &UserClaims{
		Id:    int(claims["user_id"].(float64)),
		Email: claims["email"].(string),
		Role:  claims["role"].(string),
		Exp:   int64(claims["exp"].(float64)),
	}

	return userClaims, nil
}
