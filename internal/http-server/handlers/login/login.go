package login

import (
	"context"
	"encoding/json"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sso/internal/api/response"
	"sso/internal/storage"
	"time"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
}

func New(auth storage.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getUserRequest Request
		if err := json.NewDecoder(r.Body).Decode(&getUserRequest); err != nil {
			render.JSON(w, r, response.Error(`invalid json request`))
			return
		}

		user, err := auth.GetUser(context.TODO(), getUserRequest.Email)
		if err != nil {
			render.JSON(w, r, response.Error("invalid credentials"))
			return
		}

		if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(getUserRequest.Password)); err != nil {
			render.JSON(w, r, response.Error("invalid credentials"))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "Auth",
			Value:   "authenticated",
			Expires: time.Now().Add(10000),
			Path:    "/",
		})

		render.JSON(w, r, response.OK())
	}
}
