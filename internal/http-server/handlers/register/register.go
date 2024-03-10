package register

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/badoux/checkmail"
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
		if err := validate(createUserRequest); err != nil {
			render.JSON(w, r, response.Error(err.Error()))
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

func validate(req Request) error {
	if req.Email == "" {
		return errors.New("email cannot be empty")
	}

	if err := checkmail.ValidateFormat(req.Email); err != nil {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password cannot be empty")
	}
	// You can add more validation rules here, like minimum password length
	return nil
}
