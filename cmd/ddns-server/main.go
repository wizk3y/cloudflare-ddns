package main

import (
	"cloudflare-ddns/internal/cfdns"
	"cloudflare-ddns/internal/config"
	"cloudflare-ddns/internal/domainutil"
	"cloudflare-ddns/internal/logger"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/cloudflare/cloudflare-go"
)

func init() {
	logger.InitLogger()
}

func main() {
	// Construct a new API object
	cfClient, err := cloudflare.New(config.CfApiKey, config.CfApiEmail)
	if err != nil {
		logger.Logger.Errorf("Error when create cf client, details: %v", err)
		return
	}

	// aggregate domains list
	mapDefaultTLDSubdomains := domainutil.GetMapTopLevelSubdomains(config.Domains)

	// warning if auth not required
	if !config.RequireAuth {
		logger.Logger.Warn("Authenticate is disabled, request might be use by anyone if service expose to internet")
	}

	// API handler
	http.HandleFunc("/register-ip", func(w http.ResponseWriter, r *http.Request) {
		if config.RequireAuth {
			reqUser, reqPassword, ok := r.BasicAuth()
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("basic auth not provided"))
				return
			}

			if reqUser != config.AuthUser || reqPassword != config.AuthPassword {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid username or password"))
				return
			}
		}

		ip := r.URL.Query().Get("ip")
		if ip == "" {
			ip = strings.Split(r.RemoteAddr, ":")[0]
		}
		logger.Logger.Infof("Client request register with IP: %s", ip)
		ttl, err := strconv.Atoi(r.URL.Query().Get("ttl"))
		if err != nil {
			ttl = config.TTL
		}

		var mapTLDSubdomains map[string][]string
		domains := r.URL.Query().Get("domains")
		if domains != "" {
			mapTLDSubdomains = domainutil.GetMapTopLevelSubdomains(strings.Split(domains, ","))
		} else {
			mapTLDSubdomains = mapDefaultTLDSubdomains
		}

		// get list zone
		pointedDomains := make([]string, 0)
		wg := sync.WaitGroup{}
		wg.Add(len(mapTLDSubdomains))
		for z, sds := range mapTLDSubdomains {
			go func(zoneName string, subDomains []string) {
				defer wg.Done()

				logger.Logger.Debugf("Start update DNS record `%v` for zone name %v", strings.Join(subDomains, ","), zoneName)
				err = cfdns.UpsertZoneDNSRecords(cfClient, r.Context(), zoneName, subDomains, ip, ttl)
				if err == nil {
					pointedDomains = append(pointedDomains, domainutil.BuildListDomain(zoneName, subDomains)...)
				}
			}(z, sds)
		}
		wg.Wait()

		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("Domains: %s pointed to %s", strings.Join(pointedDomains, ","), ip)
		_, err = w.Write([]byte(msg))
		if err != nil {
			logger.Logger.Errorw("Error when write response to client, details: %v", err,
				"response_msg", msg,
			)
		}
	})

	logger.Logger.Infof("Starting http server at :%d", config.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		logger.Logger.Errorw("Error when start http server, details: %v", err)
		return
	}
}
