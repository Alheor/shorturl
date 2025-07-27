// Package ip - сервис работы с клиентскими ip адресами.
//
// # Описание
//
// Проверяет разрешено ли клиенту выполнять запрос, по его ip адресу.
package ip

import (
	"net"
	"net/http"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
)

// HeaderXRealIP - имя заголовка с ip адресом
const HeaderXRealIP = `X-Real-IP`

var trustedSubnet string

// Init Подготовка обработчика IP адресов
func Init(config *config.Options) {
	trustedSubnet = config.TrustedSubnet
}

// SubnetHTTPHandler Обработчик проверки ip адреса
func SubnetHTTPHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		ipAddr := req.Header.Get(HeaderXRealIP)

		if ipAddr == `` {
			logger.Error(`empty ip address`, nil)
			resp.WriteHeader(http.StatusForbidden)

			return
		}

		if trustedSubnet == `` {
			logger.Error(`empty trusted subnet`, nil)
			resp.WriteHeader(http.StatusForbidden)

			return
		}

		if !isIPAllowed(ipAddr, trustedSubnet) {
			logger.Error(`Not allowed ip address: `+ipAddr, nil)
			resp.WriteHeader(http.StatusForbidden)

			return
		}

		f(resp, req)
	}
}

func isIPAllowed(ip string, allowedCIDRs string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	_, network, err := net.ParseCIDR(allowedCIDRs)
	if err != nil {
		return false
	}

	if network.Contains(parsedIP) {
		return true
	}

	return false
}
