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

const (
	expToken        time.Duration = time.Minute * 5    // live time of token
	expRefreshToken time.Duration = time.Hour * 24 * 7 // live time of refresh token
	UserId                        = "UserId"
	ExpToken                      = "exp" // задано библиотекой!?
)

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

func GetAccessTokenFromCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("access_token")

		if err != nil {
			handler.ServeHTTP(w, r)
			return
		}

		tknStr := c.Value                               // получаем токен из кук
		r.Header.Add("Authorization", "Bearer "+tknStr) //добавляем токен в реквест чтоб проврить в уже готово CheckJWT

		handler.ServeHTTP(w, r)
	})
}

/*
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
*/
// CreateJWTToken creates JWT token by id
func CreateJWTToken(id uint64) (string, error) {
	var err error

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", string(mySigningKey)) //TODO this should be in an env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims[UserId] = id
	atClaims[ExpToken] = time.Now().Add(expToken).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET"))) //TODO написать нормальный секретный код

	if err != nil {
		return "", err
	}

	return token, nil
}

func CreateJWTRefreshToken(id uint64) (string, error) {
	var err error

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", string(mySigningKey)) //TODO this should be in an env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims[UserId] = id
	atClaims[ExpToken] = time.Now().Add(expRefreshToken).Unix()

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
		claims, err := GetClaims(token)
		// get id from token
		if err != nil {
			return
		}

		userId := claims[UserId]

		/*if HasTokenInDB(token, 24) {
			return
		}
		*/
		// пытаемся вставить в контекст чтоб гденить еще получмить по ключу
		ctx := context.WithValue(r.Context(), userId, userId)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*func HasTokenInDB(token string, userId uint64) bool {
	conf, err := config.GetConfig()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	repoToken := base.Token()

	tokenRepo, err := repoToken.FindTokenByUserId(userId)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	if strings.Compare(token, tokenRepo.Token) == 0 {
		return true
	}

	return false
}*/

func GetClaims(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
