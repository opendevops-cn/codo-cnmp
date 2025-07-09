package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

type X509CertificateOption func(*x509.Certificate)

func WithCertificateOptionIPAddress(ips ...net.IP) X509CertificateOption {
	return func(certificate *x509.Certificate) {
		certificate.IPAddresses = ips
	}
}

func WithCertificateOptionDNSNames(names ...string) X509CertificateOption {
	return func(certificate *x509.Certificate) {
		certificate.DNSNames = names
	}
}

func GenerateCert(options ...X509CertificateOption) (cert []byte, key []byte, err error) {
	// 生成随机的序列号
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Codo, Inc."},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1年有效期
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
	for _, opt := range options {
		opt(template)
	}

	// 生成一个新的私钥
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// 创建证书
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	// 创建证书文件
	var certBs bytes.Buffer
	err = pem.Encode(&certBs, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, nil, err
	}

	// 创建私钥文件
	var keyBs bytes.Buffer
	err = pem.Encode(&keyBs, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	if err != nil {
		return nil, nil, err
	}

	return certBs.Bytes(), keyBs.Bytes(), nil
}
