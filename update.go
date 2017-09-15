package main

import (
	"crypto/x509"
	"strings"

	hierr "github.com/reconquest/hierr-go"
)

func updateCert(
	certificate string,
	privateKey string,
	suffixBac string,
	certPemData []byte,
	keyPemData []byte,
	certData *x509.Certificate,
) error {

	certDataNew, err := parseCert(certPemData)
	if err != nil {
		return hierr.Errorf(
			err, "can't parse certificate from Vault for %s",
			certData.Subject.CommonName,
		)
	}

	if certDataNew.NotAfter.Unix() <= certData.NotAfter.Unix() {
		return nil
	}
	err = checkNginx()
	if err != nil {
		return hierr.Errorf(
			err,
			"failed nginx check before update certificate %s",
			certData.Subject.CommonName,
		)
	}
	certificateBac := strings.Join([]string{certificate, suffixBac}, ".")
	privateKeyBac := strings.Join([]string{privateKey, suffixBac}, ".")

	err = copyFile(certificate, certificateBac)
	if err != nil {
		return hierr.Errorf(
			err,
			"can't copy file %s to %s",
			certificate,
			certificateBac,
		)
	}
	err = copyFile(privateKey, privateKeyBac)
	if err != nil {
		return hierr.Errorf(
			err,
			"can't copy file %s to %s",
			privateKey,
			privateKeyBac,
		)
	}
	err = writeFile(certificate, certPemData)
	if err != nil {
		return hierr.Errorf(err, "can't write file %s", certificate)
	}
	err = writeFile(privateKey, keyPemData)
	if err != nil {
		return hierr.Errorf(err, "can't write file %s", privateKey)
	}

	err = reloadNginx()
	if err != nil {
		return hierr.Errorf(
			err,
			"failed nginx reload after update certificate %s",
			certData.Subject.CommonName,
		)
	}

	return nil
}
