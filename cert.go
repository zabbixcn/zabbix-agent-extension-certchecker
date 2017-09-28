package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"

	hierr "github.com/reconquest/hierr-go"
)

func parseCert(pemData []byte) (*x509.Certificate, error) {

	block, _ := pem.Decode([]byte(pemData))

	if block == nil {
		return nil, fmt.Errorf("not found PEM data")
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("not found CERTIFICATE in PEM data")
	}

	certData, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, hierr.Errorf(err, "error from parse certificate")
	}

	certData.Subject.CommonName = strings.Replace(
		certData.Subject.CommonName,
		"*",
		"wildcard",
		1,
	)

	return certData, nil
}

func parseCertFile(certificate string) (*x509.Certificate, error) {

	pemData, err := ioutil.ReadFile(certificate)
	if err != nil {
		return nil, hierr.Errorf(
			err,
			"error read certificate %s file",
			certificate,
		)
	}

	certData, err := parseCert(pemData)
	if err != nil {
		return nil, hierr.Errorf(err, "can't parse PEM data")
	}
	return certData, nil
}

func checkCertKeyFile(certificate, privateKey string) error {

	_, err := tls.LoadX509KeyPair(certificate, privateKey)
	if err != nil {
		return hierr.Errorf(
			err,
			"error load certificate %s and key %s files",
			certificate,
			privateKey,
		)
	}

	return nil

}
