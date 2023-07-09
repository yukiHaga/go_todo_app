package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type RegisterUser struct {
	Service   RegisterUserService
	Validator *validator.Validate
}

func (userHandler *RegisterUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}

	// Struct は、構造体の公開フィールドのバリデーションを行い、 特に指定がない限り入れ子になった構造体のバリデーションも自動的に行います。
	if err := userHandler.Validator.Struct(b); err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
	}

	user, err := userHandler.Service.RegisterUser(ctx, b.Name, b.Password, b.Role)
	if err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	RespondJson(ctx, w, user, http.StatusOK)
}
