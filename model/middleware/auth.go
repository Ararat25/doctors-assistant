package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webApp/custom_errors"
	"webApp/model"
)

type AuthMiddleware struct {
	authService *model.Service
}

func NewAuthMiddleware(authService *model.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// CheckToken проверяет наличие access токена в Cookie, а также то, что он не "испортился"
func (a *AuthMiddleware) CheckToken(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("accessToken")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			apiError, _ := json.Marshal(custom_errors.ApiError{Message: "access token not found"})
			_, err = w.Write(apiError)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		userLogin, err := a.authService.VerifyUser(accessToken.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			apiError, _ := json.Marshal(custom_errors.ApiError{Message: err.Error()})
			_, err = w.Write(apiError)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		fmt.Printf("Пользователь %s - сделал запрос %s\n", userLogin, r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
