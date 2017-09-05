package main

import (
	"bytes"
	"os/exec"

	hierr "github.com/reconquest/hierr-go"
)

func nginxCheckReload() error {
	var stderr bytes.Buffer

	command := "nginx"

	path, err := exec.LookPath(command)
	if err != nil {
		return hierr.Errorf(err, "Don't find nginx bin")
	}
	cmd := exec.Command(path, "-t")
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return hierr.Errorf(stderr.String(), "Check nginx configuration failed")
	}

	cmd = exec.Command(path, "-s", "reload")
	err = cmd.Run()
	if err != nil {
		return hierr.Errorf(stderr.String(), "Reload nginx failed")
	}
	return nil
}
