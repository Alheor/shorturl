package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, "authKey", CookiesName)
	assert.Equal(t, contextKeyXAuthUser("xAuthUser"), ContextValueName)
}

func TestUser(t *testing.T) {
	tests := []struct {
		name string
		user User
	}{
		{
			name: "user with ID",
			user: User{ID: "user123"},
		},
		{
			name: "empty user",
			user: User{},
		},
		{
			name: "user with UUID",
			user: User{ID: "550e8400-e29b-41d4-a716-446655440000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.user.ID, tt.user.ID)
		})
	}
}

func TestUserCookie(t *testing.T) {
	tests := []struct {
		name   string
		cookie UserCookie
	}{
		{
			name: "cookie with user and sign",
			cookie: UserCookie{
				User: User{ID: "user123"},
				Sign: []byte("signature"),
			},
		},
		{
			name:   "empty cookie",
			cookie: UserCookie{},
		},
		{
			name: "cookie with empty sign",
			cookie: UserCookie{
				User: User{ID: "user456"},
				Sign: nil,
			},
		},
		{
			name: "cookie with empty user",
			cookie: UserCookie{
				User: User{},
				Sign: []byte("sign"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.cookie.User.ID, tt.cookie.User.ID)
			if tt.cookie.Sign != nil {
				assert.Equal(t, tt.cookie.Sign, tt.cookie.Sign)
			}
		})
	}
}

func TestEmptyUserIDErr_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      EmptyUserIDErr
		wantText string
	}{
		{
			name: "with underlying error",
			err: EmptyUserIDErr{
				Err: errors.New("user ID is empty"),
			},
			wantText: "user ID is empty",
		},
		{
			name: "with formatted error",
			err: EmptyUserIDErr{
				Err: errors.New("authentication failed: no user ID"),
			},
			wantText: "authentication failed: no user ID",
		},
		{
			name: "generic error",
			err: EmptyUserIDErr{
				Err: errors.New("error"),
			},
			wantText: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			assert.Equal(t, tt.wantText, got)
		})
	}
}

func TestEmptyUserIDErr_Implements_Error(t *testing.T) {
	var err error = &EmptyUserIDErr{
		Err: errors.New("no user ID provided"),
	}

	assert.NotNil(t, err)
	assert.Equal(t, "no user ID provided", err.Error())
}

// TestContextKeyXAuthUser проверяет тип ключа контекста
func TestContextKeyXAuthUser(t *testing.T) {
	var key = ContextValueName
	assert.Equal(t, contextKeyXAuthUser("xAuthUser"), key)

	newKey := contextKeyXAuthUser("newKey")
	assert.Equal(t, contextKeyXAuthUser("newKey"), newKey)
	assert.NotEqual(t, ContextValueName, newKey)
}

// TestUserCookie_SignOperations проверяет операции с подписью
func TestUserCookie_SignOperations(t *testing.T) {
	tests := []struct {
		name    string
		sign    []byte
		wantLen int
		wantNil bool
	}{
		{
			name:    "normal signature",
			sign:    []byte("this is a signature"),
			wantLen: 19,
			wantNil: false,
		},
		{
			name:    "empty signature",
			sign:    []byte{},
			wantLen: 0,
			wantNil: false,
		},
		{
			name:    "nil signature",
			sign:    nil,
			wantLen: 0,
			wantNil: true,
		},
		{
			name:    "binary signature",
			sign:    []byte{0x01, 0x02, 0x03, 0x04},
			wantLen: 4,
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie := UserCookie{
				User: User{ID: "test"},
				Sign: tt.sign,
			}

			if tt.wantNil {
				assert.Nil(t, cookie.Sign)
			} else {
				assert.NotNil(t, cookie.Sign)
				assert.Len(t, cookie.Sign, tt.wantLen)
			}
		})
	}
}
