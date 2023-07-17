package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/fsimic346/go-blog/model"
	"github.com/fsimic346/go-blog/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type Authenticator struct {
	JWTKey  string
	RedisDB *redis.Client
}

func (a *Authenticator) Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("jwt")
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		var claims model.UserClaim

		token, _ := jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.JWTKey), nil
		})

		if !token.Valid {
			util.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		_, err = a.RedisDB.LPos(r.Context(), "tokens", cookie.Value, redis.LPosArgs{}).Result()
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:    "jwt",
				Expires: time.Now(),
				Path:    "/",
			})
			util.RespondWithError(w, http.StatusUnauthorized, "Token manually invalidated")
			return
		}

		ctx := context.WithValue(r.Context(), "Username", claims.Username)
		newReq := r.WithContext(ctx)

		next.ServeHTTP(w, newReq)
	})
}

func (a *Authenticator) NotAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("jwt")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		token, _ := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.JWTKey), nil
		})

		if token.Valid {
			util.RespondWithError(w, http.StatusOK, "Already authenticated")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authenticator) IsAdmin(next http.HandlerFunc, userRepository model.UserRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("Username").(string)
		if !ok {
			util.RespondWithError(w, http.StatusInternalServerError, "Couldn't convert username")
			return
		}

		user, err := userRepository.GetByUsername(username)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Couldn't find user")
			return
		}

		if !user.IsAdmin {
			util.RespondWithError(w, http.StatusUnauthorized, "User not admin")
			return
		}

		ctx := context.WithValue(r.Context(), "userId", user.Id)
		newReq := r.WithContext(ctx)

		next.ServeHTTP(w, newReq)

	})
}
