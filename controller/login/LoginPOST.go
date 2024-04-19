package login

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
	"webApp/initializers"
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

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	user := model.User{
		Email:    req.PostFormValue("email"),
		Password: req.PostFormValue("password"),
	}

	userFound := model.User{}

	result := initializers.DB.Where(&model.User{Email: user.Email}).First(&userFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(res, result.Error.Error(), http.StatusNoContent)
		return
	}

	hashPassword := sha256.Sum256([]byte(user.Password))
	hashString := hex.EncodeToString(hashPassword[:])

	result = initializers.DB.Where(&model.User{Email: user.Email, Password: hashString}).First(&userFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(res, result.Error.Error(), http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub":  user.Email,
		"role": "reader",
		"exp":  time.Now().Add(time.Hour * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("TOKEN_SALT")))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonToken, err := json.Marshal(LoginResponse{AccessToken: signedToken})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.authService.AuthUser(user.Email, user.Password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(jsonToken)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
