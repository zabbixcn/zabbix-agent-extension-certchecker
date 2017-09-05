package main

import (
	"crypto/x509"
	"fmt"
	"strconv"
	"strings"

	zsend "github.com/blacked/go-zabbix"
)

func makePrefix(prefix, key string) string {
	return fmt.Sprintf(
		"%s.%s", prefix, key,
	)

}

func createCertificateMetrics(
	hostname string,
	certData *x509.Certificate,
	metrics []*zsend.Metric,
	prefix string,
	remaining int64,
) []*zsend.Metric {

	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				prefix,
				fmt.Sprintf("dnsnames.[%s]", certData.Subject.CommonName),
			),
			strings.Join(certData.DNSNames, " "),
		),
	)

	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				prefix,
				fmt.Sprintf("notbefore.[%s]", certData.Subject.CommonName),
			),
			strconv.Itoa(int(certData.NotBefore.Unix())),
		),
	)

	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				prefix,
				fmt.Sprintf("notafter.[%s]", certData.Subject.CommonName),
			),
			strconv.Itoa(int(certData.NotAfter.Unix())),
		),
	)

	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				prefix,
				fmt.Sprintf("remaining.[%s]", certData.Subject.CommonName),
			),
			strconv.Itoa(int(remaining)),
		),
	)

	return metrics
}
