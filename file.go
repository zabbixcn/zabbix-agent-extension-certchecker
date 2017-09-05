package main

import (
	"io/ioutil"
	"os"

	hierr "github.com/reconquest/hierr-go"
)

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
		return hierr.Errorf(err, "Can't read file %s", srcFilename)
	}
	err = ioutil.WriteFile(dstFilename, content, 0644)
	if err != nil {
		return hierr.Errorf(err, "Can't write file %s", dstFilename)
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
