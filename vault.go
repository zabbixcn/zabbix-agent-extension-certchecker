package main

import (
	"crypto/x509"
	"fmt"
	"strings"

	api "github.com/hashicorp/vault/api"
	hierr "github.com/reconquest/hierr-go"
)

func getCertFromVault(
	vaultAddress string,
	mountPoint string,
	tokenReadCert string,
	certData *x509.Certificate,
) ([]byte, []byte, error) {

	client, _ := api.NewClient(api.DefaultConfig())

	client.SetAddress(vaultAddress)
	client.SetToken(tokenReadCert)

	vault := client.Logical()

	resp, err := vault.Read(strings.Join(
		[]string{mountPoint, certData.Subject.CommonName}, "/"),
	)
	if err != nil || resp == nil {
		return nil, nil, hierr.Errorf(
			err,
			"can't read certificate and key for %s from Vault",
			certData.Subject.CommonName,
		)
	}

	certPemData := []byte(resp.Data["cert"].(string))
	keyPemData := []byte(resp.Data["key"].(string))

	err = checkCertKeyBlock(certPemData, keyPemData)
	if err != nil {
		return nil, nil, hierr.Errorf(err, "can't check certificate and key from Vault")
	}

	certDataNew, err := parseCert(certPemData)
	if err != nil {
		return nil, nil, hierr.Errorf(err, "can't parse certificate from Vault")
	}

	if certDataNew.NotAfter.Unix() <= certData.NotAfter.Unix() {
		return nil, nil, fmt.Errorf("certificate not renewed")
	}
	return certPemData, keyPemData, nil
}
