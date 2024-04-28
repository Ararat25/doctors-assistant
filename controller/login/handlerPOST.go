package login

import (
	"encoding/json"
	"errors"
	"net/http"
	"webApp/custom_errors"
	"webApp/model"
)

type Handler struct {
	authService *model.Service
}

func NewHandler(authService *model.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

type ResponseBody struct {
	Status bool `json:"status"`
	Id     int  `json:"id"`
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	user := model.User{
		Email:    req.PostFormValue("email"),
		Password: req.PostFormValue("password"),
	}

	id, tokens, err := h.authService.AuthUser(user.Email, user.Password)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNotFound) {
			res.WriteHeader(http.StatusNoContent)
		}
		if errors.Is(err, custom_errors.ErrIncorrectPassword) {
			res.WriteHeader(http.StatusUnauthorized)
		}
		apiError, _ := json.Marshal(custom_errors.ApiError{Message: err.Error()})

		_, err = res.Write(apiError)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	accessTokenCookie := http.Cookie{
		Name:     "accessToken",
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Path:     "/",
	}
	refreshTokenCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(res, &accessTokenCookie)
	http.SetCookie(res, &refreshTokenCookie)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	response, _ := json.Marshal(ResponseBody{Status: true, Id: id})
	_, err = res.Write(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
