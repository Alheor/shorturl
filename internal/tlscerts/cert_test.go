package tlscerts

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var certFile = baseFilePath + certFileName
var keyFile = baseFilePath + keyFileName

func TestPrepareCert(t *testing.T) {
	os.Remove(certFile)
	os.Remove(keyFile)

	// Проверка создания файлов
	certPath, keyPath, err := GenerateCert()
	require.NoError(t, err)
	assert.Equal(t, certFile, certPath)
	assert.Equal(t, keyFile, keyPath)

	// Проверка наличия файлов
	_, err = os.Stat(certPath)
	require.NoError(t, err)

	assert.False(t, os.IsNotExist(err))

	_, err = os.Stat(keyPath)
	require.NoError(t, err)

	assert.False(t, os.IsNotExist(err))

	// Проверка содержания файлов
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	require.NoError(t, err)

	assert.NotEmpty(t, cert.Certificate)

	os.Remove(certFile)
	os.Remove(keyFile)
}

func TestPrepareCert_CertificateProperties(t *testing.T) {
	os.Remove(certFile)
	os.Remove(keyFile)

	certPath, _, err := GenerateCert()
	require.NoError(t, err)

	// Читаем и парсим сертификат
	certPEM, err := os.ReadFile(certPath)
	require.NoError(t, err)

	block, _ := pem.Decode(certPEM)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	// Проверяем свойства сертификата
	var expected int64 = 1567
	assert.Equal(t, expected, cert.SerialNumber.Int64())
	assert.Equal(t, CertOrganization, cert.Subject.Organization[0])
	assert.Equal(t, CertCountry, cert.Subject.Country[0])
	assert.Equal(t, CertLocality, cert.Subject.Locality[0])

	// Проверяем IP адреса
	expectedIPs := []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback}
	if len(cert.IPAddresses) != len(expectedIPs) {
		t.Errorf("Expected %d IP addresses, got %d", len(expectedIPs), len(cert.IPAddresses))
	} else {
		for i, expectedIP := range expectedIPs {
			if !cert.IPAddresses[i].Equal(expectedIP) {
				t.Errorf("Expected IP %s at index %d, got %s", expectedIP, i, cert.IPAddresses[i])
			}
		}
	}

	if cert.NotBefore.After(time.Now()) {
		t.Error("certificate NotBefore is in the future")
	}

	expectedNotAfter := time.Now().Add(365 * 24 * time.Hour)
	timeDiff := cert.NotAfter.Sub(expectedNotAfter)
	if timeDiff < -time.Minute || timeDiff > time.Minute {
		t.Errorf("certificate NotAfter is not approximately 1 year from now")
	}

	os.Remove(certFile)
	os.Remove(keyFile)
}

func TestPrepareCert_PrivateKeyProperties(t *testing.T) {
	os.Remove(certFile)
	os.Remove(keyFile)

	_, keyPath, err := GenerateCert()
	require.NoError(t, err)

	// Читаем и парсим приватный ключ
	keyPEM, err := os.ReadFile(keyPath)
	require.NoError(t, err)

	block, _ := pem.Decode(keyPEM)
	assert.NotNil(t, block)
	assert.NotEmpty(t, block.Bytes)

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	require.NoError(t, err)

	// Проверяем что это RSA ключ правильного размера
	if rsaKey, ok := privateKey.(*rsa.PrivateKey); ok {
		if rsaKey.Size() != 512 { // 4096/8
			t.Errorf("Expected RSA key size 4096 bits, got %d bits", rsaKey.Size()*8)
		}
	} else {
		t.Error("invalid RSA private key")
	}

	os.Remove(certFile)
	os.Remove(keyFile)
}

func TestPrepareCert_FileCleanup(t *testing.T) {
	err := os.WriteFile(certFile, []byte("fake cert"), 0644)
	if err != nil {
		t.Fatalf("Failed to create fake cert file: %v", err)
	}
	err = os.WriteFile(keyFile, []byte("fake key"), 0644)
	if err != nil {
		t.Fatalf("Failed to create fake key file: %v", err)
	}

	_, _, err = GenerateCert()
	if err != nil {
		t.Fatalf("PrepareCert() returned error: %v", err)
	}

	// Проверяем что файлы были перезаписаны валидными данными
	_, err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		t.Fatalf("Failed to load certificate and key after overwrite: %v", err)
	}

	err = os.Remove(certFile)
	require.NoError(t, err)

	err = os.Remove(keyFile)
	require.NoError(t, err)
}
