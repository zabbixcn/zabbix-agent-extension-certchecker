package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	zsend "github.com/blacked/go-zabbix"
	docopt "github.com/docopt/docopt-go"
)

var version = "[manual build]"

func main() {
	usage := `zabbix-agent-extension-certchecker

Usage:
  zabbix-agent-extension-certchecker [options]
  zabbix-agent-extension-certchecker [-h --help]

Discovery options:
  --discovery               Discovery certificate file in directory.
  --path <path>             Certificate path [default: /etc/nginx/certs].
  --suffix-cert <crt>       Certificate file suffix [default: crt].
  --suffix-key <key>        Private key file suffix [default: key].

Certificate check options:
  -c --certificate <file>   Certificate file
  -k --private-key <file>   Private key file
  -z --zabbix <host>        Hostname or IP address of zabbix server
                             [default: 127.0.0.1].
  -p --port <port>          Port of zabbix server [default: 10051].
  --zabbix-prefix <prefix>  Custom prefix for key [default: certificate].

Misc options:
  --help                    Show this screen.
`

	args, _ := docopt.Parse(usage, nil, true, version, false)

	if args["--discovery"].(bool) {
		path := args["--path"].(string)
		suffixCert := args["--suffix-cert"].(string)
		suffixKey := args["--suffix-key"].(string)

		if !strings.HasSuffix(path, "/") {
			path = strings.Join([]string{path, "/"}, "")
		}
		err := discovery(path, suffixCert, suffixKey)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	certificate, ok := args["--certificate"].(string)
	if !ok {
		fmt.Println("need setup --certificate <file>")
		os.Exit(1)
	}
	privateKey, ok := args["--private-key"].(string)
	if !ok {
		fmt.Println("need setup --private-key <file>")
		os.Exit(1)
	}
	err := checkCertKeyFile(
		certificate,
		privateKey,
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	certData, err := parseCertFile(certificate)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc).Unix()

	remaining := certData.NotAfter.Unix() - now

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	zabbix := args["--zabbix"].(string)
	port, err := strconv.Atoi(args["--port"].(string))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	zabbixPrefix := args["--zabbix-prefix"].(string)

	var metrics []*zsend.Metric

	metrics = createCertificateMetrics(
		hostname,
		certData,
		metrics,
		zabbixPrefix,
		remaining,
	)

	packet := zsend.NewPacket(metrics)
	sender := zsend.NewSender(
		zabbix,
		port,
	)
	sender.Send(packet)

	fmt.Println("OK")
}
