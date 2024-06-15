package logout

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
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

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("user")

	userFound := model.User{}
	result := h.authService.Storage.Where(&model.User{Id: id}).First(&userFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)

		apiError, _ := json.Marshal(custom_errors.ErrNotFound)
		_, err := res.Write(apiError)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	_ = h.authService.Storage.Model(&model.User{}).Where(&model.User{Id: userFound.Id}).Update("accessTokenID", "")

	accessTokenCookie := http.Cookie{
		Name:     "accessToken",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
	}
	refreshTokenCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(res, &accessTokenCookie)
	http.SetCookie(res, &refreshTokenCookie)

	res.WriteHeader(http.StatusOK)
}
