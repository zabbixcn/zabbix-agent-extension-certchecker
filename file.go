package main

import (
	"io/ioutil"
	"os"
	"strings"

	hierr "github.com/reconquest/hierr-go"
)

func filterFileBySuffix(
	path,
	pattern string,
) ([]os.FileInfo, error) {

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, hierr.Errorf(err, "error read path %s", path)
	}

	var filtered []os.FileInfo

	for _, entrie := range entries {
		if strings.HasSuffix(entrie.Name(), pattern) {
			filtered = append(filtered, entrie)
		}
	}

	return filtered, nil
}

func copyFile(
	srcFilename string,
	dstFilename string,
) error {

	_, err := os.Stat(dstFilename)
	if err == nil {
		err := os.Remove(dstFilename)
		if err != nil {
			return err
		}
	}

	content, err := ioutil.ReadFile(srcFilename)
	if err != nil {
		return hierr.Errorf(err, "can't read file %s", srcFilename)
	}

	err = ioutil.WriteFile(dstFilename, content, 0644)
	if err != nil {
		return hierr.Errorf(err, "can't write file %s", dstFilename)
	}
	return nil
}

func writeFile(
	filename string,
	data []byte,
) error {

	err := os.Remove(filename)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
