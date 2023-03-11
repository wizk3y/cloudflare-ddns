package config

import (
	"flag"
	"strings"
)

var (
	Development bool
	LogDir      string

	CfApiKey   string
	CfApiEmail string
	Domains    []string
	TTL        int

	Mode string

	Port         int
	RequireAuth  bool
	AuthUser     string
	AuthPassword string
)

func init() {
	var (
		d string
	)

	flag.BoolVar(&Development, "development", false, "")
	flag.StringVar(&LogDir, "log-dir", "", "")
	flag.StringVar(&CfApiKey, "cf-api-key", "", "")
	flag.StringVar(&CfApiEmail, "cf-api-email", "", "")
	flag.StringVar(&d, "domains", "", "")
	flag.IntVar(&TTL, "ttl", 1, "")
	flag.StringVar(&Mode, "mode", "cron", "")
	flag.IntVar(&Port, "port", 8008, "")
	flag.BoolVar(&RequireAuth, "require-auth", true, "")
	flag.StringVar(&AuthUser, "auth-user", "", "")
	flag.StringVar(&AuthPassword, "auth-password", "", "")
	flag.Parse()

	if Development || Mode == "single" {
		Mode = "single"
	} else {
		Mode = "cron"
	}

	Domains = strings.Split(d, ",")
}
