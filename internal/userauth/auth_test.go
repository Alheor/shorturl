package userauth

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"testing"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestInitSuccess(t *testing.T) {
	cfg := config.Options{SignatureKey: `test_key`}
	Init(&cfg)
	assert.NotEmpty(t, signatureKey)
}

func TestGetUserSuccess(t *testing.T) {
	ctx := context.Background()
	user := models.User{}
	ctx = context.WithValue(ctx, models.ContextValueName, &user)

	assert.NotNil(t, GetUser(ctx))
}

func TestGetSignatureSuccess(t *testing.T) {
	cfg := config.Options{SignatureKey: `test_key`}
	Init(&cfg)
	assert.NotEmpty(t, signatureKey)

	sig := GetSignature(`test`)

	assert.NotNil(t, sig)
	assert.NotEmpty(t, sig)
}

func TestParseCookieInvalidLength(t *testing.T) {
	cookie := &http.Cookie{}
	cookie.Value = `1`

	userCookie, err := parseCookie(cookie)
	assert.Nil(t, err)
	assert.Nil(t, userCookie)

	cookie = &http.Cookie{}
	cookie.Value = `1234567890123456789012345678901`

	userCookie, err = parseCookie(cookie)
	assert.Nil(t, err)
	assert.Nil(t, userCookie)
}

func TestParseCookieInvalidBase64(t *testing.T) {
	cookie := &http.Cookie{}
	cookie.Value = `123456789012345678901234567890123`

	userCookie, err := parseCookie(cookie)
	assert.Nil(t, err)
	assert.Nil(t, userCookie)

	cookie = &http.Cookie{}
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(`123456789012345678901234567890123`))

	userCookie, err = parseCookie(cookie)
	assert.NotNil(t, err)
	assert.Nil(t, userCookie)
}

func TestParseCookieEmptyUser(t *testing.T) {
	cookie := &http.Cookie{}
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(`12345678901234567890123456789012`))

	userCookie, err := parseCookie(cookie)
	assert.Nil(t, userCookie)

	var myErr *models.EmptyUserIDErr
	if !errors.As(err, &myErr) {
		t.Errorf("expected: *models.EmptyUserIDErr, actual: %T", err)
	}
}

func TestParseCookieInvalidSignature(t *testing.T) {

	cookie := &http.Cookie{}
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(`123456789012345678901234567890126a30af51-b6ac-63ba-9e1c-5da06e1b610e`))

	userCookie, err := parseCookie(cookie)
	assert.Nil(t, err)
	assert.Nil(t, userCookie)
}

func TestParseCookieValidSignature(t *testing.T) {

	userId := `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`
	cookie := &http.Cookie{}
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(string(GetSignature(userId)) + userId))

	userCookie, err := parseCookie(cookie)
	assert.Nil(t, err)
	assert.NotNil(t, userCookie)
}
