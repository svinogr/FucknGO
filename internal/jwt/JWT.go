package jwt

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	jwt2 "github.com/form3tech-oss/jwt-go"
	"net/http"
	"os"
	"time"
)

var mySigningKey = []byte("SECRET")

// hadnler catch jwt token
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
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

func CreateJWT(id uint64) (string, error) {
	var err error

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", string(mySigningKey)) //TODO this should be in an env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET"))) //TODO написать нормальный секретный код

	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseJWT(token string) bool {
	/*	parseToken, err := jwt.Parse(token, func(token *jwt.Token) ([]byte, error) {
			return myLookupKey(token.Header["jdnfksdmfksd"])
		})
	*/
	/*	if err != nil {
			return false
		}

		return parseToken.Valid*/
	return false
}
