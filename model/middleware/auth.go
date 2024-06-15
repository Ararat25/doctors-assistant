package middleware

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"webApp/custom_errors"
	"webApp/initializers"
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

		id := r.URL.Query().Get("user")

		userFound := model.User{}
		result := initializers.DB.Where(&model.User{Id: id, Email: userLogin}).First(&userFound)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusUnauthorized)

			ts, err := template.ParseFiles("./view/error/notAvailable.tmpl")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = ts.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
