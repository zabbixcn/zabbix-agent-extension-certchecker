package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/reconquest/hierr-go"
)

func find(path string) ([]os.FileInfo, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, hierr.Errorf(err, "Error read path %s", path)
	}
	return entries, nil
}

func Filter(path, pattern string, filter func([]os.FileInfo, string) []os.FileInfo) ([]os.FileInfo, error) {
	entries, err := find(path)
	if err != nil {
		return nil, err
	}
	filtered := filter(entries, pattern)
	return filtered, nil
}

func hasSuffix(es []os.FileInfo, pattern string) []os.FileInfo {
	var buf []os.FileInfo
	for _, e := range es {
		if strings.HasSuffix(e.Name(), pattern) {
			buf = append(buf, e)
		}
	}
	return buf
}

func discovery(path, suffixCert, suffixKey string) error {

	discoveryData := make(map[string][]map[string]string)

	var discoveredItems []map[string]string

	suffixCert = strings.Join([]string{".", suffixCert}, "")
	suffixKey = strings.Join([]string{".", suffixKey}, "")

	filtered, err := Filter(path, suffixCert, hasSuffix)
	if err != nil {
		return err
	}
	if filtered == nil {
		return fmt.Errorf("Certificate with suffix %s not found in path %s", suffixCert, path)
	}

	for _, e := range filtered {

		certName := e.Name()

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
