// Package userauth - сервис авторизации.
//
// # Описание
//
// Авторизует пользователя по наличию специально подписанной cookie.
// Если у пользователя нет cookie или данные в ней невалидно подписаны, то выдается новая cookie.
package userauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/models"

	"github.com/google/uuid"
)

var signatureKey []byte

// Init Подготовка сервиса к работе.
func Init(config *config.Options) {
	signatureKey = []byte(config.SignatureKey)
}

// GetUser get user from context.
func GetUser(ctx context.Context) *models.User {
	authUser := ctx.Value(models.ContextValueName)
	if authUser != nil {
		return authUser.(*models.User)
	}

	return nil
}

// AuthHTTPHandler обработчик авторизации пользователя.
func AuthHTTPHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		var userCookie *models.UserCookie
		var err error

		cookie, _ := req.Cookie(models.CookiesName)
		if cookie != nil {
			userCookie, err = parseCookie(cookie)
		}

		if err != nil {
			var myErr *models.EmptyUserIDErr
			if errors.As(err, &myErr) {
				f(resp, req)
				return
			}
		}

		if userCookie == nil {

			userCookie = getUserCookie()
			cookieValue := string(userCookie.Sign) + userCookie.User.ID

			http.SetCookie(resp,
				&http.Cookie{
					Name:  models.CookiesName,
					Value: base64.StdEncoding.EncodeToString([]byte(cookieValue)),
				},
			)
		}

		ctxWithUser := context.WithValue(req.Context(), models.ContextValueName, &userCookie.User)
		f(resp, req.Clone(ctxWithUser))
	}
}

// GetSignature Получение подписи.
func GetSignature(uuid string) []byte {
	h := hmac.New(sha256.New, signatureKey)
	h.Write([]byte(models.CookiesName))
	h.Write([]byte(uuid))

	return h.Sum(nil)
}

func parseCookie(cookie *http.Cookie) (userCookie *models.UserCookie, error error) {

	if len(cookie.Value) < sha256.Size {
		return nil, nil
	}

	cookieValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, nil
	}

	signature := cookieValue[:sha256.Size]
	value, err := uuid.Parse(string(cookieValue)[sha256.Size:])

	if err != nil {
		return nil, &models.EmptyUserIDErr{Err: err}
	}

	expectedSignature := GetSignature(value.String())

	if !hmac.Equal(signature, expectedSignature) {
		return nil, nil
	}

	return &models.UserCookie{User: models.User{ID: value.String()}, Sign: signature}, nil
}

func getUserCookie() *models.UserCookie {
	userID := generateUserUUID()
	return &models.UserCookie{User: models.User{ID: userID}, Sign: GetSignature(userID)}
}

func generateUserUUID() string {
	return uuid.New().String()
}

// ParseCookieToken парсит строковый токен (для gRPC)
func ParseCookieToken(token string) (*models.User, error) {
	if len(token) == 0 {
		return nil, errors.New("empty token")
	}

	// Если токен закодирован в base64, декодируем его
	cookieValue, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		// Возможно, токен уже декодирован, попробуем использовать как есть
		cookieValue = []byte(token)
	}

	if len(cookieValue) < sha256.Size {
		return nil, errors.New("token too short")
	}

	signature := cookieValue[:sha256.Size]
	userIDBytes := cookieValue[sha256.Size:]

	userUUID, err := uuid.Parse(string(userIDBytes))
	if err != nil {
		return nil, &models.EmptyUserIDErr{Err: err}
	}

	expectedSignature := GetSignature(userUUID.String())

	if !hmac.Equal(signature, expectedSignature) {
		return nil, errors.New("invalid signature")
	}

	return &models.User{ID: userUUID.String()}, nil
}
