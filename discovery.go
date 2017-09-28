package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func discovery(path, suffixCert, suffixKey string) error {

	discoveryData := make(map[string][]map[string]string)

	var discoveredItems []map[string]string

	suffixCert = strings.Join([]string{".", suffixCert}, "")
	suffixKey = strings.Join([]string{".", suffixKey}, "")

	filtered, err := filterFileBySuffix(path, suffixCert)
	if err != nil {
		return err
	}
	if filtered == nil {
		return fmt.Errorf(
			"certificate with suffix %s not found in path %s",
			suffixCert,
			path,
		)
	}

	for _, entrie := range filtered {
		certName := entrie.Name()
		keyName := strings.Replace(certName, suffixCert, suffixKey, 1)
		certificate := strings.Join([]string{path, certName}, "")
		privateKey := strings.Join([]string{path, keyName}, "")

		err = checkCertKeyFile(certificate, privateKey)
		if err != nil {
			continue
		}

		certData, err := parseCertFile(certificate)
		if err != nil {
			continue
		}

		discoveredItem := make(map[string]string)
		discoveredItem["{#COMMONNAME}"] = certData.Subject.CommonName
		discoveredItem["{#CERTIFICATE}"] = certificate
		discoveredItem["{#PRIVATEKEY}"] = privateKey
		discoveredItems = append(discoveredItems, discoveredItem)
	}
	discoveryData["data"] = discoveredItems

	out, err := json.Marshal(discoveryData)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", out)
	return nil
}
