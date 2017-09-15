package main

import (
	"os/exec"

	hierr "github.com/reconquest/hierr-go"
)

func checkNginx() error {
	cmd := exec.Command("/usr/bin/nginx", "-t")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return hierr.Errorf(stdoutStderr, "check nginx configuration failed")
	}
	return nil
}
func reloadNginx() error {
	err := checkNginx()
	if err != nil {
		return hierr.Errorf(err, "check nginx configuration failed")
	}

	cmd := exec.Command("/usr/bin/nginx", "-s", "reload")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return hierr.Errorf(stdoutStderr, "reload nginx failed")
	}
	return nil
}
