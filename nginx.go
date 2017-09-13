package main

import (
	"bytes"
	"os/exec"

	hierr "github.com/reconquest/hierr-go"
)

func nginxCheck() error {
	var stderr bytes.Buffer

	cmd := exec.Command("/usr/bin/nginx", "-t")
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return hierr.Errorf(stderr.String(), "check nginx configuration failed")
	}
	return nil
}
func nginxReload() error {
	var stderr bytes.Buffer

	err := nginxCheck()
	if err != nil {
		return hierr.Errorf(stderr.String(), "check nginx configuration failed")
	}

	cmd := exec.Command("/usr/bin/nginx", "-s", "reload")
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return hierr.Errorf(stderr.String(), "reload nginx failed")
	}
	return nil
}
