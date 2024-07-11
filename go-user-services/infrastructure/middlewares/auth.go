package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	logger "github.com/wlrudi19/library-management-app/go-user-services/utils/log"
	"github.com/wlrudi19/library-management-app/go-user-services/utils/response"
)

const (
	ContextKeyUserEmail = "userEmail"
	ContextKeyUserId    = "userId"
	CorrID              = "X-Correlation-ID"
	SecretKey           = "x-secret-key"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenString := request.Header.Get("Authorization")

		if tokenString == "" {
			resp := response.CustomBuilder(http.StatusUnauthorized, false, nil, "You are not authorized to access this resource")
			resp.Send(writer)
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			resp := response.CustomBuilder(http.StatusUnauthorized, false, nil, "Token invalid")
			resp.Send(writer)
			return
		}

		email, ok := claims["Email"].(string)
		if !ok {
			resp := response.CustomBuilder(http.StatusBadRequest, false, nil, "Token invalid")
			resp.Send(writer)
			return
		}

		userIdFloat, ok := claims["Id"].(float64)
		if !ok {
			resp := response.CustomBuilder(http.StatusUnauthorized, false, nil, "Token invalid")
			resp.Send(writer)
			return
		}
		userId := int(userIdFloat)

		ctx := context.WithValue(request.Context(), ContextKeyUserEmail, email)
		ctx = context.WithValue(request.Context(), ContextKeyUserId, userId)
		ctx = context.WithValue(request.Context(), CorrID, uuid.NewString())
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GenerateAccessToken(userId int, email string) (string, error) {
	//generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":    userId,
		"Email": email,
		"exp":   time.Now().Add(time.Minute * 60).Unix(), //time expired 1 hour
	})

	//sign token
	secretKey := []byte("x-library-rudi")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		logger.GetRequestLogEntry(context.Background(), "GenerateAccessToken", email).Error(err)
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (map[string]interface{}, error) {
	secretKey := []byte("x-library-rudi")

	//parsing & validate hashing method
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims) //return is bool
	if !ok || !token.Valid {
		return nil, errors.New("token invalid")
	}

	expTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if expTime.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

func InternalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		secretKey := request.Header.Get("X-Secret-Key")

		if secretKey == "" || secretKey != string(SecretKey) {
			resp := response.CustomBuilder(http.StatusUnauthorized, false, nil, "You are not authorized to access this resource")
			resp.Send(writer)
			return
		}

		ctx := context.WithValue(request.Context(), CorrID, uuid.NewString())
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
