package ip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Options
		expectedSubnet string
	}{
		{
			name: "valid config",
			config: &config.Options{
				TrustedSubnet: "192.168.1.0/24",
			},
			expectedSubnet: "192.168.1.0/24",
		},
		{
			name: "empty trusted subnet",
			config: &config.Options{
				TrustedSubnet: "",
			},
			expectedSubnet: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.config)
			if trustedSubnet != tt.expectedSubnet {
				t.Errorf("Init() trustedSubnet = %v, want %v", trustedSubnet, tt.expectedSubnet)
			}
		})
	}
}

func TestIsIPAllowed(t *testing.T) {
	tests := []struct {
		name         string
		ip           string
		allowedCIDRs string
		want         bool
	}{
		{
			name:         "valid IP in subnet",
			ip:           "192.168.1.100",
			allowedCIDRs: "192.168.1.0/24",
			want:         true,
		},
		{
			name:         "valid IP not in subnet",
			ip:           "10.0.0.1",
			allowedCIDRs: "192.168.1.0/24",
			want:         false,
		},
		{
			name:         "invalid IP address",
			ip:           "not-an-ip",
			allowedCIDRs: "192.168.1.0/24",
			want:         false,
		},
		{
			name:         "invalid CIDR",
			ip:           "192.168.1.100",
			allowedCIDRs: "invalid-cidr",
			want:         false,
		},
		{
			name:         "IPv6 address in subnet",
			ip:           "2001:db8::1",
			allowedCIDRs: "2001:db8::/32",
			want:         true,
		},
		{
			name:         "IPv6 address not in subnet",
			ip:           "2001:db9::1",
			allowedCIDRs: "2001:db8::/32",
			want:         false,
		},
		{
			name:         "localhost IP",
			ip:           "127.0.0.1",
			allowedCIDRs: "127.0.0.0/8",
			want:         true,
		},
		{
			name:         "single IP CIDR",
			ip:           "10.0.0.1",
			allowedCIDRs: "10.0.0.1/32",
			want:         true,
		},
		{
			name:         "edge of subnet",
			ip:           "192.168.1.255",
			allowedCIDRs: "192.168.1.0/24",
			want:         true,
		},
		{
			name:         "outside edge of subnet",
			ip:           "192.168.2.1",
			allowedCIDRs: "192.168.1.0/24",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIPAllowed(tt.ip, tt.allowedCIDRs); got != tt.want {
				t.Errorf("isIPAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubnetHTTPHandler(t *testing.T) {

	err := logger.Init(nil)
	assert.NoError(t, err)

	// Настраиваем trustedSubnet для тестов
	originalTrustedSubnet := trustedSubnet
	defer func() {
		trustedSubnet = originalTrustedSubnet
	}()

	// Создаем тестовый handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	tests := []struct {
		name               string
		trustedSubnet      string
		xRealIP            string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "valid IP in trusted subnet",
			trustedSubnet:      "192.168.1.0/24",
			xRealIP:            "192.168.1.100",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "success",
		},
		{
			name:               "valid IP not in trusted subnet",
			trustedSubnet:      "192.168.1.0/24",
			xRealIP:            "10.0.0.1",
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       "",
		},
		{
			name:               "empty X-Real-IP header",
			trustedSubnet:      "192.168.1.0/24",
			xRealIP:            "",
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       "",
		},
		{
			name:               "empty trusted subnet",
			trustedSubnet:      "",
			xRealIP:            "192.168.1.100",
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       "",
		},
		{
			name:               "invalid IP address",
			trustedSubnet:      "192.168.1.0/24",
			xRealIP:            "not-an-ip",
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       "",
		},
		{
			name:               "IPv6 address in trusted subnet",
			trustedSubnet:      "2001:db8::/32",
			xRealIP:            "2001:db8::1",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "success",
		},
		{
			name:               "localhost access",
			trustedSubnet:      "127.0.0.0/8",
			xRealIP:            "127.0.0.1",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем trustedSubnet для каждого теста
			trustedSubnet = tt.trustedSubnet

			// Создаем запрос
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			// Создаем ResponseRecorder для записи ответа
			rr := httptest.NewRecorder()

			// Создаем обработчик с middleware
			handler := SubnetHTTPHandler(testHandler)

			// Выполняем запрос
			handler.ServeHTTP(rr, req)

			// Проверяем статус код
			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatusCode)
			}

			// Проверяем тело ответа
			if body := rr.Body.String(); body != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					body, tt.expectedBody)
			}
		})
	}
}

func TestSubnetHTTPHandlerIntegration(t *testing.T) {
	// Интеграционный тест, проверяющий полный цикл работы
	conf := &config.Options{
		TrustedSubnet: "192.168.0.0/16",
	}
	Init(conf)

	err := logger.Init(nil)
	assert.NoError(t, err)

	var handlerCalled bool
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := SubnetHTTPHandler(testHandler)

	t.Run("allowed IP calls handler", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Real-IP", "192.168.1.1")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if !handlerCalled {
			t.Error("expected handler to be called for allowed IP")
		}
		if rr.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", rr.Code)
		}
	})

	t.Run("forbidden IP does not call handler", func(t *testing.T) {
		handlerCalled = false
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Real-IP", "10.0.0.1")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if handlerCalled {
			t.Error("expected handler NOT to be called for forbidden IP")
		}
		if rr.Code != http.StatusForbidden {
			t.Errorf("expected status Forbidden, got %d", rr.Code)
		}
	})
}

// Бенчмарк для функции isIPAllowed
func BenchmarkIsIPAllowed(b *testing.B) {
	ip := "192.168.1.100"
	cidr := "192.168.1.0/24"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isIPAllowed(ip, cidr)
	}
}

// Бенчмарк для SubnetHTTPHandler
func BenchmarkSubnetHTTPHandler(b *testing.B) {
	trustedSubnet = "192.168.1.0/24"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := SubnetHTTPHandler(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Real-IP", "192.168.1.100")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
