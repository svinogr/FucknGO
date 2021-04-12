package jwt

import (
	"FucknGO/db/repo"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"context"
	"errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	jwt2 "github.com/form3tech-oss/jwt-go"
	"net/http"
	"os"
	"strings"
	"time"
)

var MySigningKey = []byte("SECRET") //TODO поменяит ключ

const (
	expToken        time.Duration = time.Second * 5    // live time of token
	expRefreshToken time.Duration = time.Hour * 24 * 7 // live time of refresh token
	UserId                        = "UserId"
	ExpToken                      = "exp" // задано стандартом
	Claims                        = "claims"
	Refresh                       = "refresh"
	User                          = "user"
)

// hadnler catch jwt token
var JwtVerifMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt2.Token) (interface{}, error) {
		return MySigningKey, nil
	},
	UserProperty:        "",
	ErrorHandler:        nil,
	CredentialsOptional: false,
	Extractor:           nil,
	Debug:               false,
	EnableAuthOnOptions: false,
	SigningMethod:       jwt.SigningMethodHS256,
})

/*type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
*/

func CheckTokensInCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refresh := false
		cR, err := r.Cookie(model.RefreshTokenName) // получаем Refresh

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized) // ошибка идем логинится так как смысла без рефреш вообще нет
			return
		}

		_, err = r.Cookie(model.AccessTokenName) // получаем Access

		if err != nil {
			refresh = true // если ошибка то возможен только рефреш

			_, err := jwt2.Parse(cR.Value, func(token *jwt2.Token) (interface{}, error) {
				return MySigningKey, nil // парсим рефреш
			})

			if err != nil { // если рефреш ошибка то идем и логинимся
				http.Error(w, err.Error(), http.StatusUnauthorized) // ошибка идем логинится так как смысла без рефреш вообще нет
				return
			}

			context.WithValue(r.Context(), Refresh, refresh) // идем рефрешится и далее
			handler.ServeHTTP(w, r)

		}

		ctx := context.WithValue(r.Context(), Refresh, refresh) // все ок. идем далее
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AccessOrRefresh(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		refresh := c.Value(Refresh).(bool)

		if refresh {
			sessionOld, ok := GetValidSessionByCookie(r) // получаем  сессию по куки

			if !ok {
				http.Error(w, errors.New("session expired").Error(), http.StatusUnauthorized)
				return
			}

			// создаем новые токены
			newAccessToken, err := CreateJWTToken(sessionOld.UserId, model.AccessTokenName) // не боимся нила так как нил невозможен

			if err != nil {
				log.NewLog().Fatal(err)
			}

			newRefreshToken, err := CreateJWTToken(sessionOld.UserId, model.RefreshTokenName)

			if err != nil {
				log.NewLog().Fatal(err)
			}

			CreateNewSessionForToken(sessionOld, newRefreshToken) // создаем новую сессию и удаляем старую

			SetCookieWithToken(&w, newAccessToken)
			SetCookieWithToken(&w, newRefreshToken)
		}

		err := SetUserToRequest(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})

}

func GetValidSessionByCookie(r *http.Request) (*repo.SessionModelRepo, bool) {
	cookie, err := r.Cookie(model.RefreshTokenName) // получаем куки по ключу рефреш

	if err != nil {
		return nil, false
	}

	byCookie, err := FindSessionByCookie(*cookie)

	if err != nil {
		return nil, false
	}

	ok := ValidSession(byCookie, r) // проверяем сессию на валид

	if !ok {
		return nil, false
	}

	return byCookie, true
}

func CreateNewSessionForToken(session *repo.SessionModelRepo, tokenModel model.TokenModel) (*repo.SessionModelRepo, error) {
	db := repo.NewDataBaseWithConfig()
	sessionsRepo := db.Sessions()
	sessionsRepo.DeleteSessionByUserId(session.UserId) // удаляем старую сессию

	session.RefreshToken = tokenModel.Value
	session.ExpireIn = time.Now().Add(repo.Exp_session)
	session.CreatedAt = time.Now()

	return sessionsRepo.CreateSession(session)
}

func ValidSession(session *repo.SessionModelRepo, r *http.Request) bool {
	cookie, err := r.Cookie(model.RefreshTokenName) // получаем куки по ключу рефреш

	if err != nil {
		return false
	}
	refreshTokenFromSessionRepo := session.RefreshToken

	if cookie.Value != refreshTokenFromSessionRepo { // проверяем ключи
		return false
	}

	in := session.ExpireIn

	if in.Before(time.Now()) {
		return false
	}

	if session.UserAgent != r.UserAgent() {
		return false
	}
	// данная проверка не работает почему то
	/*	if session.Ip != r.RemoteAddr {
			return false
		}
	*/
	return true

}

func SetUserToRequest(r *http.Request) (err error) {
	cookie, err := r.Cookie(model.RefreshTokenName) // получаем куки по ключу рефреш

	user, err := FindUserByCookie(cookie)

	c := context.WithValue(r.Context(), User, user)
	r.WithContext(c)

	return err
}

func FindUserByCookie(cookie *http.Cookie) (*repo.UserModelRepo, error) {
	db := repo.NewDataBaseWithConfig()
	userRepo := db.User()

	tokenStr := cookie.Value

	claims := jwt2.MapClaims{}

	_, err := jwt2.ParseWithClaims(tokenStr, claims, func(token *jwt2.Token) (interface{}, error) {
		return MySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	userId := uint64(claims[UserId].(float64))

	return userRepo.FindUserById(userId)
}

func FindSessionByCookie(cookie http.Cookie) (*repo.SessionModelRepo, error) {
	tokenStr := cookie.Value

	claims := jwt2.MapClaims{}

	_, err := jwt2.ParseWithClaims(tokenStr, claims, func(token *jwt2.Token) (interface{}, error) {
		return MySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	userId := uint64(claims[UserId].(float64))

	db := repo.NewDataBaseWithConfig()
	sessionsRepo := db.Sessions()

	return sessionsRepo.FindSessionByUserId(userId)
}

/*func GetAccessTokenFromCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(model.AccessTokenName) // получаем куку

		if err != nil {
			http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
			return
		}

		tokenValue := c.Value // получаем значение токена из куки

		claims := jwt2.MapClaims{}

		token, err := jwt2.ParseWithClaims(tokenValue, claims, func(token *jwt2.Token) (interface{}, error) {
			return MySigningKey, nil
		})

		if err != nil {
			fmt.Print(err.Error())
			http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
			return
		}

		if token.Valid {
			ctx := context.WithValue(r.Context(), Claims, claims)
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// если что то не так с токеном проверяем рефреш
		c, err = r.Cookie(model.RefreshTokenName) // получаем куку

		if err != nil {
			http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
			return
		}

		tokenValue = c.Value // получаем значение токена из куки

		claims = jwt2.MapClaims{}

		token, err = jwt2.ParseWithClaims(tokenValue, claims, func(token *jwt2.Token) (interface{}, error) {
			return MySigningKey, nil
		})

		if err != nil {
			http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
			return
		}

		if token.Valid {
			userId := uint64(claims[UserId].(float64))
			base := repo.NewDataBaseWithConfig()
			sessionRepo := base.Sessions()

			session, err := sessionRepo.GetSessionForUserIdIfIs(userId)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			refreshTokenFromSessionRepo := session.RefreshToken

			//если токен соответсвует токену в базе сессий
			if refreshTokenFromSessionRepo != token.Raw {
				if _, err := sessionRepo.DeleteSessionByUserId(userId); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} // удаляем сессию
				// создаем новые токены
				accessToken, _ := CreateJWTToken(userId)
				refreshToken, _ := CreateJwtRefreshToken(userId)

				newSession := repo.SessionModelRepo{
					UserId:       userId,
					RefreshToken: refreshToken.Value,
					UserAgent:    r.UserAgent(),
					Fingerprint:  "",
					Ip:           r.RemoteAddr,
					ExpireIn:     time.Now().Add(repo.Exp_session),
					CreatedAt:    time.Now(),
				}

				sessionRepo.CreateSession(&newSession)
				// добавляем httpOnly Cookie
				SetCookieWithToken(&w, accessToken)
				SetCookieWithToken(&w, refreshToken)
				ctx := context.WithValue(r.Context(), Claims, claims)
				handler.ServeHTTP(w, r.WithContext(ctx))
				return
			}

		}

		http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
	})
}
*/ /*func GetAccessTokenFromCookie2(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cR, err := r.Cookie(model.RefreshTokenName) // получаем куку

		if err != nil {
			http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
			return
		}

		c, err := r.Cookie(model.AccessTokenName) // получаем куку

		if err != nil {
			http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
			return
		}

		tokenValue := c.Value // получаем значение токена из куки

		claims := jwt2.MapClaims{}

		token, err := jwt2.ParseWithClaims(tokenValue, claims, func(token *jwt2.Token) (interface{}, error) {
			return MySigningKey, nil
		})

		if err != nil {
			fmt.Print(err.Error())
			http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
			return
		}

		if token.Valid {
			ctx := context.WithValue(r.Context(), Claims, claims)
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// если что то не так с токеном проверяем рефреш
		//-----------------------

		tokenValue = c.Value // получаем значение токена из куки

		claims = jwt2.MapClaims{}

		token, err = jwt2.ParseWithClaims(tokenValue, claims, func(token *jwt2.Token) (interface{}, error) {
			return MySigningKey, nil
		})

		if err != nil {
			http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
			return
		}

		if token.Valid {
			userId := uint64(claims[UserId].(float64))
			base := repo.NewDataBaseWithConfig()
			sessionRepo := base.Sessions()

			session, err := sessionRepo.GetSessionForUserIdIfIs(userId)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			refreshTokenFromSessionRepo := session.RefreshToken

			//если токен соответсвует токену в базе сессий
			if refreshTokenFromSessionRepo != token.Raw {
				if _, err := sessionRepo.DeleteSessionByUserId(userId); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} // удаляем сессию
				// создаем новые токены
				accessToken, _ := CreateJWTToken(userId)
				refreshToken, _ := CreateJwtRefreshToken(userId)

				newSession := repo.SessionModelRepo{
					UserId:       userId,
					RefreshToken: refreshToken.Value,
					UserAgent:    r.UserAgent(),
					Fingerprint:  "",
					Ip:           r.RemoteAddr,
					ExpireIn:     time.Now().Add(repo.Exp_session),
					CreatedAt:    time.Now(),
				}

				sessionRepo.CreateSession(&newSession)
				// добавляем httpOnly Cookie
				SetCookieWithToken(&w, accessToken)
				SetCookieWithToken(&w, refreshToken)
				ctx := context.WithValue(r.Context(), Claims, claims)
				handler.ServeHTTP(w, r.WithContext(ctx))
				return
			}

		}

		http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
	})
}
*/
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
			return MySigningKey, nil
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
func CreateJWTToken(id uint64, nameToken string) (model.TokenModel, error) {
	var err error

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", string(MySigningKey)) //TODO this should be in an env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims[UserId] = id

	var expTime time.Time

	if nameToken == model.AccessTokenName {
		expTime = time.Now().Add(expToken)
		atClaims[ExpToken] = time.Now().Add(expToken).Unix()
	}

	if nameToken == model.RefreshTokenName {
		expTime = time.Now().Add(expRefreshToken)
		atClaims[ExpToken] = time.Now().Add(expRefreshToken).Unix()
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	tokenValue, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET"))) //TODO написать нормальный секретный код

	if err != nil {
		return model.TokenModel{}, err
	}

	token := model.TokenModel{
		Name:    model.AccessTokenName,
		Value:   tokenValue,
		ExpTime: expTime,
	}

	return token, nil
}

/*func CreateJwtRefreshToken(id uint64) (model.TokenModel, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", string(MySigningKey)) //TODO this should be in an env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims[UserId] = id
	expTime := time.Now().Add(expRefreshToken)
	atClaims[ExpToken] = time.Now().Add(expRefreshToken).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	tokenValue, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET"))) //TODO написать нормальный секретный код

	if err != nil {
		return model.TokenModel{}, err
	}

	token := model.TokenModel{
		Name:    model.RefreshTokenName,
		Value:   tokenValue,
		ExpTime: expTime,
	}

	return token, nil
}*/

// CreateCookieWithToken creates  cookie with token by id
func CreateCookieWithToken(name string, value string, expTime time.Time) http.Cookie {
	cookie := http.Cookie{
		Name:       name,
		Value:      value,
		Path:       "/",
		Domain:     "",
		Expires:    expTime,
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   true, // attention
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}

	return cookie
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
		//TODO возможно стоит сюда проверку userId вставить
		ctx := context.WithValue(r.Context(), UserId, userId)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return MySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func SetCookieWithToken(w *http.ResponseWriter, token model.TokenModel) {
	cookieWithToken := CreateCookieWithToken(token.Name, token.Value, token.ExpTime)
	http.SetCookie(*w, &cookieWithToken)
}

func GetUserIdFromContext(r *http.Request) (repo.UserModelRepo, error) {
	value := r.Context().Value(User)

	if value == nil {
		return repo.UserModelRepo{}, errors.New("Not found user")
	}

	return value.(repo.UserModelRepo), nil
}

func DeleteCookie(w *http.ResponseWriter) {
	cookie := http.Cookie{
		Name:       "",
		Value:      "",
		Path:       "",
		Domain:     "",
		Expires:    time.Unix(0, 0),
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}

	http.SetCookie(*w, &cookie)
}
