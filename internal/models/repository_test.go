package models

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchEl(t *testing.T) {
	tests := []struct {
		name string
		el   BatchEl
	}{
		{
			name: "full element",
			el: BatchEl{
				CorrelationID: "123",
				OriginalURL:   "https://example.com",
				ShortURL:      "http://short.url/abc",
			},
		},
		{
			name: "empty element",
			el:   BatchEl{},
		},
		{
			name: "partial element",
			el: BatchEl{
				OriginalURL: "https://example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.el.CorrelationID, tt.el.CorrelationID)
			assert.Equal(t, tt.el.OriginalURL, tt.el.OriginalURL)
			assert.Equal(t, tt.el.ShortURL, tt.el.ShortURL)
		})
	}
}

func TestUniqueErr_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      UniqueErr
		wantText string
	}{
		{
			name: "with underlying error",
			err: UniqueErr{
				ShortKey: "abc123",
				Err:      errors.New("duplicate key"),
			},
			wantText: "duplicate key",
		},
		{
			name: "with formatted error",
			err: UniqueErr{
				ShortKey: "xyz789",
				Err:      errors.New("key xyz789 already exists"),
			},
			wantText: "key xyz789 already exists",
		},
		{
			name: "empty short key",
			err: UniqueErr{
				ShortKey: "",
				Err:      errors.New("error"),
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

func TestUniqueErr_Implements_Error(t *testing.T) {
	var err error = &UniqueErr{
		ShortKey: "test",
		Err:      errors.New("test error"),
	}

	assert.NotNil(t, err)
	assert.Equal(t, "test error", err.Error())
}

func TestHistoryNotFoundErr(t *testing.T) {
	tests := []struct {
		name string
		err  HistoryNotFoundErr
	}{
		{
			name: "empty error",
			err:  HistoryNotFoundErr{},
		},
		{
			name: "with embedded error",
			err: HistoryNotFoundErr{
				error: errors.New("no history found"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.error != nil {
				assert.NotNil(t, tt.err.error)
			}
		})
	}
}

func TestHistoryEl_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		el      HistoryEl
		want    string
		wantErr bool
	}{
		{
			name: "valid history element",
			el: HistoryEl{
				OriginalURL: "https://example.com",
				ShortURL:    "http://short.url/abc123",
			},
			want: `{"original_url":"https://example.com","short_url":"http://short.url/abc123"}`,
		},
		{
			name: "empty fields",
			el:   HistoryEl{},
			want: `{"original_url":"","short_url":""}`,
		},
		{
			name: "only original URL",
			el: HistoryEl{
				OriginalURL: "https://example.com",
			},
			want: `{"original_url":"https://example.com","short_url":""}`,
		},
		{
			name: "only short URL",
			el: HistoryEl{
				ShortURL: "http://short.url/xyz",
			},
			want: `{"original_url":"","short_url":"http://short.url/xyz"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.el)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestHistoryEl_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    HistoryEl
		wantErr bool
	}{
		{
			name: "valid JSON",
			data: `{"original_url":"https://example.com","short_url":"http://short.url/abc123"}`,
			want: HistoryEl{
				OriginalURL: "https://example.com",
				ShortURL:    "http://short.url/abc123",
			},
		},
		{
			name: "empty JSON",
			data: `{}`,
			want: HistoryEl{},
		},
		{
			name: "partial JSON",
			data: `{"original_url":"https://example.com"}`,
			want: HistoryEl{
				OriginalURL: "https://example.com",
			},
		},
		{
			name:    "invalid JSON",
			data:    `{invalid}`,
			wantErr: true,
		},
		{
			name: "extra fields",
			data: `{"original_url":"https://example.com","short_url":"http://short.url/abc","extra":"field"}`,
			want: HistoryEl{
				OriginalURL: "https://example.com",
				ShortURL:    "http://short.url/abc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got HistoryEl
			err := json.Unmarshal([]byte(tt.data), &got)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestUniqueErr_ShortKey проверяет доступность поля ShortKey
func TestUniqueErr_ShortKey(t *testing.T) {
	tests := []struct {
		name     string
		shortKey string
	}{
		{
			name:     "with short key",
			shortKey: "abc123",
		},
		{
			name:     "empty short key",
			shortKey: "",
		},
		{
			name:     "special characters",
			shortKey: "!@#$%^",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UniqueErr{
				ShortKey: tt.shortKey,
				Err:      errors.New("test"),
			}
			assert.Equal(t, tt.shortKey, err.ShortKey)
		})
	}
}
