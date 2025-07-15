package tlscerts

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareCustomCert_SuccessfulLoad(t *testing.T) {
	os.Remove(certFileName)
	os.Remove(keyFileName)

	certPEM, keyPEM, err := generateTestCertificate()
	require.NoError(t, err)

	certBase64 := base64.StdEncoding.EncodeToString([]byte(certPEM))
	keyBase64 := base64.StdEncoding.EncodeToString([]byte(keyPEM))

	certFile, keyFile, err := prepareCustomCert(certBase64, keyBase64)
	require.NoError(t, err)

	assert.Equal(t, certFile, certFileName)
	assert.Equal(t, keyFile, keyFileName)

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		t.Errorf("Certificate file was not created")
	}

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		t.Errorf("Key file was not created")
	}

	certContent, err := os.ReadFile(certFile)
	require.NoError(t, err)
	assert.Equal(t, string(certContent), certPEM)

	keyContent, err := os.ReadFile(keyFile)
	require.NoError(t, err)
	assert.Equal(t, string(keyContent), keyPEM)

	os.Remove(certFileName)
	os.Remove(keyFileName)
}

func TestPrepareCustomCert_InvalidCertificateBase64(t *testing.T) {
	os.Remove(certFileName)
	os.Remove(keyFileName)

	invalidCertBase64 := "invalid-base64-certificate!!!"
	validKeyBase64 := base64.StdEncoding.EncodeToString([]byte("valid-key-content"))

	_, _, err := prepareCustomCert(invalidCertBase64, validKeyBase64)
	require.Error(t, err)

	os.Remove(certFileName)
	os.Remove(keyFileName)
}

func TestPrepareCustomCert_InvalidKeyBase64(t *testing.T) {
	os.Remove(certFileName)
	os.Remove(keyFileName)

	validCertBase64 := base64.StdEncoding.EncodeToString([]byte("valid-certificate-content"))
	invalidKeyBase64 := "invalid-base64-key!!!"

	_, _, err := prepareCustomCert(validCertBase64, invalidKeyBase64)
	require.Error(t, err)

	os.Remove(certFileName)
	os.Remove(keyFileName)
}

func TestPrepareCustomCert_EmptyBase64Values(t *testing.T) {
	os.Remove(certFileName)
	os.Remove(keyFileName)

	_, _, err := prepareCustomCert("", "")
	require.NoError(t, err)

	if _, err := os.Stat(certFileName); os.IsNotExist(err) {
		t.Errorf("Certificate file was not created")
	}

	if _, err := os.Stat(keyFileName); os.IsNotExist(err) {
		t.Errorf("Key file was not created")
	}

	certContent, err := os.ReadFile(certFileName)
	require.NoError(t, err)
	assert.True(t, len(certContent) == 0)

	keyContent, err := os.ReadFile(keyFileName)
	require.NoError(t, err)
	assert.True(t, len(keyContent) == 0)

	os.Remove(certFileName)
	os.Remove(keyFileName)
}

func TestGetCert_WithValidCustomCertificates(t *testing.T) {
	os.Remove(certFileName)
	os.Remove(keyFileName)

	certPEM, keyPEM, err := generateTestCertificate()
	require.NoError(t, err)

	certBase64 := base64.StdEncoding.EncodeToString([]byte(certPEM))
	keyBase64 := base64.StdEncoding.EncodeToString([]byte(keyPEM))

	certFile, keyFile, err := GetCert(certBase64, keyBase64)
	require.NoError(t, err)
	assert.Equal(t, certFile, certFileName)
	assert.Equal(t, keyFile, keyFileName)

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		t.Errorf("Certificate file was not created")
	}

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		t.Errorf("Key file was not created")
	}

	certContent, err := os.ReadFile(certFile)
	require.NoError(t, err)

	assert.Equal(t, string(certContent), certPEM)

	os.Remove(certFileName)
	os.Remove(keyFileName)
}

func TestGetCert_WithoutCustomCertificates(t *testing.T) {
	os.Remove(certFileName)
	os.Remove(keyFileName)

	certFile, keyFile, err := GetCert("", "")
	require.NoError(t, err)

	assert.Equal(t, certFile, certFileName)
	assert.Equal(t, keyFile, keyFileName)

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		t.Errorf("Certificate file was not created")
	}

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		t.Errorf("Key file was not created")
	}

	os.Remove(certFileName)
	os.Remove(keyFileName)
}

func TestGetCert_WithPartialCustomData(t *testing.T) {
	os.Remove(certFileName)
	os.Remove(keyFileName)

	testCases := []struct {
		name     string
		certData string
		keyData  string
	}{
		{"Only certificate provided", "some-cert-data", ""},
		{"Only key provided", "", "some-key-data"},
		{"Certificate with empty key", base64.StdEncoding.EncodeToString([]byte("cert")), ""},
		{"Key with empty certificate", "", base64.StdEncoding.EncodeToString([]byte("key"))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Remove(certFileName)
			os.Remove(keyFileName)

			certFile, keyFile, err := GetCert(tc.certData, tc.keyData)
			require.NoError(t, err)

			assert.Equal(t, certFile, certFileName)
			assert.Equal(t, keyFile, keyFileName)

			if _, err := os.Stat(certFile); os.IsNotExist(err) {
				t.Errorf("Certificate file was not created")
			}

			if _, err := os.Stat(keyFile); os.IsNotExist(err) {
				t.Errorf("Key file was not created")
			}
		})
	}

	os.Remove(certFileName)
	os.Remove(keyFileName)
}

func TestPrepareCustomCert_FileCreationError(t *testing.T) {

	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	certBase64 := base64.StdEncoding.EncodeToString([]byte(certPEM))
	keyBase64 := base64.StdEncoding.EncodeToString([]byte(keyPEM))

	err = os.Mkdir(certFileName, 0755)
	require.NoError(t, err)

	defer os.Remove(certFileName)

	_, _, err = prepareCustomCert(certBase64, keyBase64)
	require.Error(t, err)
}

func TestPrepareCustomCert_LargeBase64Data(t *testing.T) {
	os.Remove(certFileName)
	os.Remove(keyFileName)

	largeCertData := make([]byte, 10000) // 10KB данных
	largeKeyData := make([]byte, 8000)   // 8KB данных

	for i := range largeCertData {
		largeCertData[i] = byte(i % 256)
	}

	for i := range largeKeyData {
		largeKeyData[i] = byte((i + 100) % 256)
	}

	certBase64 := base64.StdEncoding.EncodeToString(largeCertData)
	keyBase64 := base64.StdEncoding.EncodeToString(largeKeyData)

	certFile, keyFile, err := prepareCustomCert(certBase64, keyBase64)
	require.NoError(t, err)

	assert.Equal(t, certFile, certFileName)
	assert.Equal(t, keyFile, keyFileName)

	certContent, err := os.ReadFile(certFile)
	require.NoError(t, err)
	assert.True(t, len(certContent) == len(largeCertData))

	keyContent, err := os.ReadFile(keyFile)
	require.NoError(t, err)
	assert.True(t, len(keyContent) == len(largeKeyData))

	os.Remove(certFileName)
	os.Remove(keyFileName)
}

// generateTestCertificate создает тестовый сертификат и ключ для тестирования
func generateTestCertificate() (certPEM, keyPEM string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Test Company"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	certPEMBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	privDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	keyPEMBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})

	return string(certPEMBytes), string(keyPEMBytes), nil
}
