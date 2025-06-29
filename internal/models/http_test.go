package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIRequest_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		req     APIRequest
		want    string
		wantErr bool
	}{
		{
			name: "valid request",
			req:  APIRequest{URL: "https://example.com"},
			want: `{"url":"https://example.com"}`,
		},
		{
			name: "empty URL",
			req:  APIRequest{URL: ""},
			want: `{"url":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestAPIRequest_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    APIRequest
		wantErr bool
	}{
		{
			name: "valid JSON",
			data: `{"url":"https://example.com"}`,
			want: APIRequest{URL: "https://example.com"},
		},
		{
			name: "empty JSON",
			data: `{}`,
			want: APIRequest{URL: ""},
		},
		{
			name: "extra fields",
			data: `{"url":"https://example.com","extra":"field"}`,
			want: APIRequest{URL: "https://example.com"},
		},
		{
			name:    "invalid JSON",
			data:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got APIRequest
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

func TestAPIResponse_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		resp    APIResponse
		want    string
		wantErr bool
	}{
		{
			name: "success response",
			resp: APIResponse{
				Result:     "http://short.url/abc123",
				StatusCode: 200,
			},
			want: `{"result":"http://short.url/abc123"}`,
		},
		{
			name: "error response",
			resp: APIResponse{
				Error:      "invalid URL",
				StatusCode: 400,
			},
			want: `{"error":"invalid URL"}`,
		},
		{
			name: "empty fields omitted",
			resp: APIResponse{
				StatusCode: 200,
			},
			want: `{}`,
		},
		{
			name: "both result and error",
			resp: APIResponse{
				Result:     "result",
				Error:      "error",
				StatusCode: 500,
			},
			want: `{"result":"result","error":"error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.resp)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestAPIBatchRequestEl_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		req     APIBatchRequestEl
		want    string
		wantErr bool
	}{
		{
			name: "valid batch request element",
			req: APIBatchRequestEl{
				CorrelationID: "123",
				OriginalURL:   "https://example.com",
			},
			want: `{"correlation_id":"123","original_url":"https://example.com"}`,
		},
		{
			name: "empty fields",
			req:  APIBatchRequestEl{},
			want: `{"correlation_id":"","original_url":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestAPIBatchResponseEl_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		resp    APIBatchResponseEl
		want    string
		wantErr bool
	}{
		{
			name: "valid batch response element",
			resp: APIBatchResponseEl{
				CorrelationID: "123",
				ShortURL:      "http://short.url/abc123",
			},
			want: `{"correlation_id":"123","short_url":"http://short.url/abc123"}`,
		},
		{
			name: "empty fields",
			resp: APIBatchResponseEl{},
			want: `{"correlation_id":"","short_url":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.resp)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestAPIBatchRequestEl_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    APIBatchRequestEl
		wantErr bool
	}{
		{
			name: "valid JSON",
			data: `{"correlation_id":"123","original_url":"https://example.com"}`,
			want: APIBatchRequestEl{
				CorrelationID: "123",
				OriginalURL:   "https://example.com",
			},
		},
		{
			name: "snake_case JSON fields",
			data: `{"correlation_id":"456","original_url":"https://test.com"}`,
			want: APIBatchRequestEl{
				CorrelationID: "456",
				OriginalURL:   "https://test.com",
			},
		},
		{
			name:    "invalid JSON",
			data:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got APIBatchRequestEl
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

// TestAPIResponse_StatusCodeNotSerialized проверяет, что StatusCode не сериализуется в JSON
func TestAPIResponse_StatusCodeNotSerialized(t *testing.T) {
	resp := APIResponse{
		Result:     "test",
		StatusCode: 200,
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	// Проверяем, что в JSON нет поля status_code или statusCode
	assert.NotContains(t, string(data), "status_code")
	assert.NotContains(t, string(data), "statusCode")
	assert.NotContains(t, string(data), "StatusCode")
}
