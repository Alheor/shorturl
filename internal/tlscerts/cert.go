// Package tlscerts - сервис сертификации HTTPS
//
// # Описание
//
// Формирует самоподписанный TLS сертификат с работы веб-сервера по протоколу HTTPS
package tlscerts

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

// CertOrganization - организация
var CertOrganization = `Short url`

// CertCountry - Код страны
var CertCountry = `RU`

// CertLocality - город
var CertLocality = `Moscow`

var certFileName = `cert.pem`
var keyFileName = `key.pem`

// GetCert получение сертификата TLS, либо из конфига, либо самоподписанного
func GetCert(c string, k string) (cert string, key string, err error) {

	if c != "" && k != "" {
		return prepareCustomCert(c, k)
	}

	return prepareCert()
}

// prepareCert формирование самоподписанного сертификата TLS
func prepareCert() (cert string, key string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	crt := &x509.Certificate{
		SerialNumber: big.NewInt(1567),
		Subject: pkix.Name{
			Organization: []string{CertOrganization},
			Country:      []string{CertCountry},
			Locality:     []string{CertLocality},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // 1 year
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, crt, crt, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	if err != nil {
		return "", "", err
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	if err != nil {
		return "", "", err
	}

	certOut, err := os.Create(certFileName)
	if err != nil {
		return "", "", err
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return "", "", err
	}

	keyOut, err := os.Create(keyFileName)
	if err != nil {
		return "", "", err
	}
	defer keyOut.Close()

	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER}); err != nil {
		return "", "", err
	}

	return certFileName, keyFileName, nil
}

// prepareCustomCert обработка пользовательских TLS сертификата и ключа в формате base64
func prepareCustomCert(certBase64, keyBase64 string) (cert string, key string, err error) {
	// Декодируем сертификат из base64
	certBytes, err := base64.StdEncoding.DecodeString(certBase64)
	if err != nil {
		return "", "", err
	}

	// Декодируем ключ из base64
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return "", "", err
	}

	// Создаем временные файлы для сертификата и ключа
	certOut, err := os.Create(certFileName)
	if err != nil {
		return "", "", err
	}
	defer certOut.Close()

	_, err = certOut.Write(certBytes)
	if err != nil {
		return "", "", err
	}

	keyOut, err := os.Create(keyFileName)
	if err != nil {
		return "", "", err
	}
	defer keyOut.Close()

	_, err = keyOut.Write(keyBytes)
	if err != nil {
		return "", "", err
	}

	return certFileName, keyFileName, nil
}
