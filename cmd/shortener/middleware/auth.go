package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/context"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

//var userId string

const TokenExp = time.Hour * 3
const SecretKey = "my-256-bit-secret"

func BuildJWTString() (string, string, error) {

	id := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: id,
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, id, nil
}

func GetUserID(tokenString string) (string, string) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err == nil {
		return tokenString, claims.UserID
	}

	tokenString, userID, _ := BuildJWTString()
	return tokenString, userID
}

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		cookies := r.Cookies()
		tokenWithUserID, err := r.Cookie("user_id")

		//Кука есть, смотрим на наличие user id. Если его нет, можно не заниматься валидацией и сразу вернуть 401
		if len(cookies) != 0 && errors.Is(err, http.ErrNoCookie) {
			ow.WriteHeader(http.StatusUnauthorized)
			return
		}

		//Если куки нет, выдаем новую, пропускаем запрос

		if tokenWithUserID == nil {
			tokenWithUserID, userID := GetUserID("")

			cookie := http.Cookie{Name: "user_id", Value: tokenWithUserID}
			http.SetCookie(w, &cookie)
			context.Set(r, "userID", userID)
			h.ServeHTTP(ow, r)
		} else {
			//Кука есть, user id есть. Проверяем подпись. Если не сошлась, выдаем новую и пропускаем. Если сошлась, то просто пропускаем

			tokenString, userID := GetUserID(tokenWithUserID.Value)
			cookie := http.Cookie{Name: "user_id", Value: tokenString}
			http.SetCookie(w, &cookie)
			context.Set(r, "userID", userID)
			h.ServeHTTP(ow, r)
		}
	}
}
