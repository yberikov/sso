package editUser

import (
	"context"
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
	"sso/internal/api/response"
	"sso/internal/storage"
)

type Request struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	DateBirth string `json:"date-birth"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(auth storage.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var editUserRequest Request
		if err := json.NewDecoder(r.Body).Decode(&editUserRequest); err != nil {
			render.JSON(w, r, response.Error(`invalid json request`))
			return
		}
		if err := auth.EditUser(context.TODO(), editUserRequest.Email, editUserRequest.Name, editUserRequest.Telephone, editUserRequest.DateBirth); err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}
