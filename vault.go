package main

import (
	"crypto/x509"
	"strings"

	api "github.com/hashicorp/vault/api"
	hierr "github.com/reconquest/hierr-go"
)

func updateFromVault(
	certificate string,
	privateKey string,
	suffixBac string,
	vaultAddress string,
	mountPoint string,
	tokenReadCert string,
	certData *x509.Certificate,
) error {

	clientConfig := api.DefaultConfig()
	client, _ := api.NewClient(clientConfig)

	client.SetAddress(vaultAddress)
	client.SetToken(tokenReadCert)

	vault := client.Logical()

	resp, err := vault.Read(strings.Join(
		[]string{mountPoint, certData.Subject.CommonName}, "/"),
	)
	if err != nil || resp == nil {
		return hierr.Errorf(
			err,
			"Can't read certificate and key for %s from Vault",
			certData.Subject.CommonName,
		)
	}

	certPemData := []byte(resp.Data["cert"].(string))
	keyPemData := []byte(resp.Data["key"].(string))

	err = checkCertKeyBlock(certPemData, keyPemData)
	if err != nil {
		return hierr.Errorf(err, "Can't check certificate and key from Vault")
	}

	certDataNew, err := parseCert(certPemData)
	if err != nil {
		return hierr.Errorf(err, "Can't parse certificate from Vault")
	}

	if certDataNew.NotAfter.Unix() > certData.NotAfter.Unix() {

		errNginx := nginxCheckReload()
		if errNginx != nil {
			return hierr.Errorf(errNginx,
				"Failed nginx check before update certificate",
			)
		}

		suffixBac := "backup"

		certificateBac := strings.Join([]string{certificate, suffixBac}, ".")
		privateKeyBac := strings.Join([]string{privateKey, suffixBac}, ".")

		err = copyFile(certificate, certificateBac)
		if err != nil {
			return hierr.Errorf(err, "Can't copy file %s to %s",
				certificate, certificateBac,
			)
		}
		err = copyFile(privateKey, privateKeyBac)
		if err != nil {
			return hierr.Errorf(err, "Can't copy file %s to %s",
				privateKey, privateKeyBac,
			)
		}

		err = writeFile(certificate, certPemData)
		if err != nil {
			return hierr.Errorf(err, "Can't write file %s", certificate)
		}
		err = writeFile(privateKey, keyPemData)
		if err != nil {
			return hierr.Errorf(err, "Can't write file %s", privateKey)
		}

		errNginx = nginxCheckReload()
		if errNginx != nil {

			err = copyFile(certificateBac, certificate)
			if err != nil {
				return hierr.Errorf(err, "Can't copy file %s to %s",
					certificateBac, certificate,
				)
			}
			err = copyFile(privateKeyBac, privateKey)
			if err != nil {
				return hierr.Errorf(err, "Can't copy file %s to %s",
					privateKeyBac, privateKey,
				)
			}

			return errNginx
		}
	}

	return nil

}
