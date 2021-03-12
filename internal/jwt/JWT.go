package jwt

import (
	"context"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	jwt2 "github.com/form3tech-oss/jwt-go"
	"net/http"
	"os"
	"strings"
	"time"
)

var mySigningKey = []byte("SECRET")

const EXP time.Duration = 300 // live time of token
const USER_ID = "userId"

// hadnler catch jwt token
var JwtVerifMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt2.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	UserProperty:        "",
	ErrorHandler:        nil,
	CredentialsOptional: false,
	Extractor:           nil,
	Debug:               false,
	EnableAuthOnOptions: false,
	SigningMethod:       jwt.SigningMethodHS256,
})

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CookieMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// если нет авторизованого токена то го хом
		c, err := r.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := c.Value
		claims := Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			handler.ServeHTTP(w, r)
		}

	})
}

// CreateJWT creates JWT token by id
func CreateJWT(id uint64) (string, error) {
	var err error

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", string(mySigningKey)) //TODO this should be in an env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims[USER_ID] = id
	atClaims["EXP"] = time.Now().Add(time.Minute * EXP).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET"))) //TODO написать нормальный секретный код

	if err != nil {
		return "", err
	}

	return token, nil
}

// ParseJWT middleware parses token to get id user
func ParseJWT(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get header
		header := r.Header.Get("Authorization")
		// получаем токен отбрасываем Bearer
		token := strings.Split(header, " ")[1]
		// default function for parsing token into claim
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

		if err != nil {
			return
		}
		// get id from token
		userId := claims[USER_ID]
		// пытаемся вставить в контекст чтоб гденить еще получмить по ключу
		ctx := context.WithValue(r.Context(), USER_ID, userId)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
