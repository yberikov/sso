package register

import (
	"context"
	"encoding/json"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sso/internal/api/response"
	"sso/internal/storage"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(auth storage.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var createUserRequest Request
		if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
			render.JSON(w, r, response.Error(`invalid json request`))
			return
		}
		if createUserRequest.Email == "" {
			render.JSON(w, r, response.Error("Empty"))
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), bcrypt.MinCost)
		if err != nil {
			render.JSON(w, r, response.Error("Internal error"))
			return
		}
		if _, err = auth.SaveUser(context.TODO(), createUserRequest.Email, hash); err != nil {
			render.JSON(w, r, response.Error("db request"))
			return
		}

		render.JSON(w, r, response.OK())
	}
}
