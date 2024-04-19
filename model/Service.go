package model

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"webApp/custom_errors"
	"webApp/initializers"
	"webApp/model/auth/entity"
)

type Service struct {
	passwordSalt []byte
	tokenSalt    []byte

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	userStorage     *gorm.DB
}

func NewAuthService(passwordSalt, tokenSalt []byte) *Service {
	return &Service{
		passwordSalt:    passwordSalt,
		tokenSalt:       tokenSalt,
		accessTokenTTL:  time.Minute,
		refreshTokenTTL: 24 * time.Hour,
		userStorage:     initializers.DB,
	}
}

// RegisterUser регистрирует нового пользователя
func (s *Service) RegisterUser(email, password, lastName, firstName, middleName, specialty string) error {
	passwordHash := s.hashPassword(password)

	user := NewUser(email, passwordHash, lastName, firstName, middleName, specialty)

	result := s.userStorage.Create(&user)
	if result.Error != nil {
		return custom_errors.ErrUserAlreadyExists
	}

	return nil
}

// AuthUser генерирует refresh и access токены для пользователя после входа в систему
func (s *Service) AuthUser(email, password string) (entity.Tokens, error) {
	user := &User{
		Email:    email,
		Password: password,
	}

	userFound := User{}

	result := s.userStorage.Where(User{Email: user.Email}).First(&userFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entity.Tokens{}, custom_errors.ErrNotFound
	}

	isPasswordCorrect := s.doPasswordsMatch(s.hashPassword(password), user.Password)
	if !isPasswordCorrect {
		return entity.Tokens{}, custom_errors.ErrIncorrectPassword
	}

	tokens, err := s.generateTokens(email)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

// VerifyUser верифицирует пользователя по access токену
func (s *Service) VerifyUser(token string) (string, error) {
	claims := &entity.AuthClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect method")
		}

		return s.tokenSalt, nil
	})
	if err != nil || !parsedToken.Valid {
		return "", fmt.Errorf("incorrect token: %v", err)
	}

	return claims.Email, nil
}

// RefreshToken обновляет токены пользователя
func (s *Service) RefreshToken(token string) (entity.Tokens, error) {
	claims := &entity.RefreshTokenClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("incorrect method")
		}

		return s.tokenSalt, nil
	})

	if err != nil || !parsedToken.Valid {
		return entity.Tokens{}, fmt.Errorf("incorrect refresh token: %v", err)
	}

	// поиск токена в хранилище
	userFound := User{}

	result := initializers.DB.Where(&User{Email: claims.Email, AccessTokenID: claims.AccessTokenID}).First(&userFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entity.Tokens{}, custom_errors.ErrNotFound
	}

	// валидация прошла успешно, можем генерить новую пару
	tokens, err := s.generateTokens(claims.Email)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

// generateTokens генерирует Access и Refresh токены
func (s *Service) generateTokens(login string) (entity.Tokens, error) {
	accessTokenID := uuid.NewString()
	accessToken, err := s.generateAccessToken(login)
	if err != nil {
		return entity.Tokens{}, err
	}

	refreshToken, err := s.generateRefreshToken(login, accessTokenID)
	if err != nil {
		return entity.Tokens{}, err
	}

	res := s.userStorage.Model(&User{}).Where(User{Email: login}).Update("accessTokenID", accessTokenID)
	if res.Error != nil {
		return entity.Tokens{}, res.Error
	}

	return entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// генерирует Access токен
func (s *Service) generateAccessToken(email string) (string, error) {
	now := time.Now()
	claims := entity.AuthClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)), // TTL - time to live
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.tokenSalt)
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return signedToken, nil
}

// generateRefreshToken генерирует Refresh токен
func (s *Service) generateRefreshToken(email, accessTokenID string) (string, error) {
	now := time.Now()
	claims := entity.RefreshTokenClaims{
		Email:         email,
		AccessTokenID: accessTokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.tokenSalt)
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return signedToken, nil
}

// hashPassword хэширует строку
func (s *Service) hashPassword(password string) string {
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()

	passwordBytes = append(passwordBytes, s.passwordSalt...)
	sha512Hasher.Write(passwordBytes)

	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

// doPasswordsMatch сравнивает хеш паролей
func (s *Service) doPasswordsMatch(hashedPassword, currPassword string) bool {
	return hashedPassword == s.hashPassword(currPassword)
}