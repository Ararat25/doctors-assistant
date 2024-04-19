package register

import (
	"net/http"
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
	user := model.User{
		Email:      req.PostFormValue("email"),
		Password:   req.PostFormValue("password"),
		LastName:   req.PostFormValue("lastName"),
		FirstName:  req.PostFormValue("firstName"),
		MiddleName: req.PostFormValue("middleName"),
		Specialty:  req.PostFormValue("specialty"),
	}

	err := h.authService.RegisterUser(user.Email, user.Password, user.LastName, user.FirstName, user.MiddleName, user.Specialty)

	if err != nil {
		http.Error(res, err.Error(), http.StatusConflict)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}
