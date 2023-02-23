package utils

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
)

var queryList = []string{
	"https://ip.3322.net",
	"https://myip.ipip.net",
}

var ipTester = regexp.MustCompile(`(\d+\.){3}\d+`)

func GetIP() (ip string, err error) {
	var resp *http.Response
	var bs []byte
	var errMsg = "failed to get self ip"
	for _, url := range queryList {
		L := log.WithField("url", url)

		resp, err = http.Get(url)
		if err != nil {
			L.WithError(err).Error(errMsg)
			continue
		}

		if bs, err = io.ReadAll(resp.Body); err != nil {
			L.WithError(err).Error(errMsg)

			continue
		}
		defer resp.Body.Close()

		bodyStr := string(bs)
		if ip = ipTester.FindString(bodyStr); ip == "" {
			err = fmt.Errorf("ip not found, original resp: %s", bodyStr)
			L.WithError(err).Error(errMsg)

			continue
		}
		if ip != "" {
			err = nil
			L.WithField("ip", ip).Debugf("found ip")
			break
		}
	}
	return
}
