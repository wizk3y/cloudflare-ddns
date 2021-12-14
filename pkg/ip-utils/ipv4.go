package ip_utils

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var ipCheckSvcs = []string{
	"https://api.ipify.org",
	"https://checkip.amazonaws.com",
	"https://v4.ident.me/",
	"https://ifconfig.me/ip",
	"https://ipv4.icanhazip.com/",
}

func GetCurrentIPv4() string {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, s := range ipCheckSvcs {
		resp, err := httpClient.Get(s)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}

			return strings.TrimSpace(string(bodyBytes))
		}
	}

	return ""
}
