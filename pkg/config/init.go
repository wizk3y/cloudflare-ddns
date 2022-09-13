package config

import (
	"flag"
	"strings"
)

var (
	Development bool
	Mode        string

	CfApiKey   string
	CfApiEmail string
	Domains    []string
)

func init() {
	var (
		d string
	)

	flag.BoolVar(&Development, "development", false, "")
	flag.StringVar(&Mode, "mode", "cron", "")
	flag.StringVar(&CfApiKey, "cf-api-key", "", "")
	flag.StringVar(&CfApiEmail, "cf-api-email", "", "")
	flag.StringVar(&d, "domains", "", "")
	flag.Parse()

	if Development && Mode != "single" {
		Mode = "single"
	} else {
		Mode = "cron"
	}

	Domains = strings.Split(d, ",")
}
