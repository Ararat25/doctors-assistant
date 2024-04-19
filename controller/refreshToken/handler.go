package refreshToken

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
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		apiError, _ := json.Marshal(custom_errors.ApiError{Message: "refresh token not found"})

		_, err = w.Write(apiError)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	tokens, err := h.authService.RefreshToken(refreshToken.Value)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNotFound) {
			apiError, _ := json.Marshal(custom_errors.ApiError{Message: "token not found in storage"})
			w.WriteHeader(http.StatusNotFound)

			_, err = w.Write(apiError)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}
		apiError, _ := json.Marshal(custom_errors.ApiError{Message: "cannot refresh token"})
		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write(apiError)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	response, _ := json.Marshal(ResponseBody{Status: true})
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
