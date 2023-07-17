package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fsimic346/go-blog/model"
	"github.com/fsimic346/go-blog/util"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type UserController struct {
	UserService model.UserService
	JWTKey      string
	RedisDB     *redis.Client
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	user, err := uc.UserService.GetById(userId)
	if err != nil {
		log.Printf("Error while fetching user: %v", err)
		util.RespondWithError(w, http.StatusNotFound, "Couldn't find user")
		return
	}

	util.RespondWithJSON(w, 200, user)
}

func (uc *UserController) AddUser(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"is_admin"`
	}
	var params reqParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		// log.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode params")
		return
	}

	user, err := uc.UserService.Add(params.Username, params.Password, params.IsAdmin)

	if err != nil {
		// log.Println(err)
		util.RespondWithError(w, http.StatusBadRequest, "User already exists")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, user)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Username string
		Password string
	}

	var params reqParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		// log.Println(err)
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode params")
		return
	}

	err = uc.UserService.Login(params.Username, params.Password)
	if err != nil {
		// log.Print(err)
		util.RespondWithError(w, http.StatusBadRequest, "Invalid login credentials")
		return
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		Username:         params.Username,
	})

	tokenString, err := token.SignedString([]byte(uc.JWTKey))

	err = uc.RedisDB.LPush(r.Context(), "tokens", tokenString).Err()
	if err != nil {
		panic(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	})
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("jwt")
	tokenString := cookie.Value

	err := uc.RedisDB.LRem(r.Context(), "tokens", 1, tokenString)
	if err != nil {
		log.Print(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Now(),
		Path:    "/",
	})
}
