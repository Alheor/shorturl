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

func Init(config *config.Options) {
	signatureKey = []byte(config.SignatureKey)
}

// GetUser get user from context
func GetUser(ctx context.Context) *models.User {
	authUser := ctx.Value(models.ContextValueName)
	if authUser != nil {
		return authUser.(*models.User)
	}

	return nil
}

func AuthHTTPHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		var userCookie *models.UserCookie
		var err error

		for _, cookie := range req.Cookies() {
			if cookie.Name == models.CookiesName {
				userCookie, err = parseCookie(cookie)
				break
			}
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
		println(userCookie.User.ID)

		ctxWithUser := context.WithValue(req.Context(), models.ContextValueName, &userCookie.User)
		f(resp, req.Clone(ctxWithUser))
	}
}

func parseCookie(cookie *http.Cookie) (userCookie *models.UserCookie, error error) {

	if len(cookie.Value) < sha256.Size {
		return nil, &models.EmptyUserIDErr{Err: nil}
	}

	cookieValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, &models.EmptyUserIDErr{Err: err}
	}

	signature := cookieValue[:sha256.Size]
	value, err := uuid.Parse(string(cookieValue)[sha256.Size:])

	if err != nil {
		return nil, &models.EmptyUserIDErr{Err: err}
	}

	expectedSignature := GetSignature(value.String())

	if !hmac.Equal(signature, expectedSignature) {
		return nil, errors.New(`invalid signature`)
	}

	return &models.UserCookie{User: models.User{ID: value.String()}, Sign: signature}, nil
}

func getUserCookie() *models.UserCookie {
	userID := generateUserUUID()
	return &models.UserCookie{User: models.User{ID: userID}, Sign: GetSignature(userID)}
}

func GetSignature(uuid string) []byte {
	h := hmac.New(sha256.New, signatureKey)
	h.Write([]byte(models.CookiesName))
	h.Write([]byte(uuid))

	return h.Sum(nil)
}

func generateUserUUID() string {
	return uuid.New().String()
}
