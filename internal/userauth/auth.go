package userauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/google/uuid"
	"net/http"
)

const CookiesName = `authKey`
const ContextValueName = `xAuthUser`

type User struct {
	ID string
}

type UserCookie struct {
	user User
	sign []byte
}

type EmptyUserIDErr struct {
	Err error
}

func (e *EmptyUserIDErr) Error() string {
	return e.Err.Error()
}

// GetUserFromContext get user from context
func GetUserFromContext(ctx context.Context) *User {

	authUser := ctx.Value(ContextValueName)

	switch authUser.(type) {
	default:
		return nil
	case *User:
		return authUser.(*User)
	}
}

func WithUserAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var userCookie *UserCookie
		var err error

		for _, cookie := range r.Cookies() {
			if cookie.Name == CookiesName {
				userCookie, err = parseCookie(cookie)
				break
			}
		}

		if err != nil {
			var myErr *EmptyUserIDErr
			if errors.As(err, &myErr) {
				f(w, r)
				return
			}
		}

		if userCookie == nil {

			userCookie = getUserCookie()
			cookieValue := string(userCookie.sign) + userCookie.user.ID

			http.SetCookie(w,
				&http.Cookie{
					Name:  CookiesName,
					Value: base64.StdEncoding.EncodeToString([]byte(cookieValue)),
				},
			)
		}

		ctxWithUser := context.WithValue(r.Context(), ContextValueName, &userCookie.user)

		f(w, r.Clone(ctxWithUser))
	}
}

func parseCookie(cookie *http.Cookie) (userCookie *UserCookie, error error) {

	cookieValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, errors.New(`invalid signature`)
	}

	signature := cookieValue[:sha256.Size]
	value, err := uuid.Parse(string(cookieValue)[sha256.Size:])

	if err != nil {
		return nil, &EmptyUserIDErr{err}
	}

	expectedSignature := GetSignature(value.String())

	if !hmac.Equal(signature, expectedSignature) {
		return nil, errors.New(`invalid signature`)
	}

	return &UserCookie{user: User{value.String()}, sign: signature}, nil
}

func getUserCookie() *UserCookie {
	userID := generateUserUuid()
	return &UserCookie{user: User{userID}, sign: GetSignature(userID)}
}

func GetSignature(uuid string) []byte {
	h := hmac.New(sha256.New, []byte(config.Options.SignatureKey))
	h.Write([]byte(CookiesName))
	h.Write([]byte(uuid))

	return h.Sum(nil)
}

func generateUserUuid() string {
	return uuid.New().String()
}
