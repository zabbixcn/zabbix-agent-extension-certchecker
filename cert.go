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

	var block *pem.Block

	block, _ = pem.Decode([]byte(pemData))

	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("Error decode pem data")
	}

	certData, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, hierr.Errorf(err, "Error from parse certificate")
	}

	certData.Subject.CommonName = strings.Replace(certData.Subject.CommonName,
		"*", "wildcard", 1)

	return certData, nil
}

func parseCertFile(certificate string) (*x509.Certificate, error) {

	pemData, err := ioutil.ReadFile(certificate)
	if err != nil {
		return nil, hierr.Errorf(err, "Error read certificate %s file",
			certificate,
		)
	}

	certData, err := parseCert(pemData)
	if err != nil {
		return nil, err
	}
	return certData, nil
}

func checkCertKeyFile(certificate, privateKey string) error {

	_, err := tls.LoadX509KeyPair(certificate, privateKey)
	if err != nil {
		return hierr.Errorf(err, "Error load certificate %s and key %s files",
			certificate, privateKey,
		)
	}

	return nil

}

func checkCertKeyBlock(certPemData, keyPemData []byte) error {

	_, err := tls.X509KeyPair(certPemData, keyPemData)
	if err != nil {
		return hierr.Errorf(err, "Error load certificate and key blocks")
	}

	return nil

}
