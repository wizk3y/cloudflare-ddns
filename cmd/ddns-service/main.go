package main

import (
	"cloudflare-ddns/internal/cfdns"
	"cloudflare-ddns/internal/config"
	"cloudflare-ddns/internal/domainutil"
	"cloudflare-ddns/internal/iputil"
	"cloudflare-ddns/internal/logger"
	"context"
	"strings"
	"sync"

	"github.com/cloudflare/cloudflare-go"
	"github.com/robfig/cron"
)

var (
	// this like an state when running in cron mode
	currentIP string
)

func init() {
	logger.InitLogger()
}

func main() {
	if config.Mode == "single" {
		pointToCurrentIP()
	} else {
		c := cron.New()
		err := c.AddFunc("0 */5 * * * *", pointToCurrentIP)
		if err != nil {
			logger.Logger.Fatalf("Error when add cron function, details: %v", err)
		}

		logger.Logger.Infof("Cron function added successfully.")

		c.Start()
		select {}
	}
}

func pointToCurrentIP() {
	var ipv4 = iputil.GetCurrentIPv4()
	// var ipv6 = iputil.GetCurrentIPv6()

	logger.Logger.Infof("Current IP v4: %v", ipv4)

	// check is IP changed
	if ipv4 == currentIP {
		logger.Logger.Debugf("The current IP of the machine has not changed.")
		return
	}

	// Construct a new API object
	cfClient, err := cloudflare.New(config.CfApiKey, config.CfApiEmail)
	if err != nil {
		logger.Logger.Errorf("Error when create cf client, details: %v", err)
		return
	}

	// Most API calls require a Context
	ctx := context.Background()

	// aggregate domains list
	mapTLDSubdomains := domainutil.GetMapTopLevelSubdomains(config.Domains)

	// get list zone
	wg := sync.WaitGroup{}
	wg.Add(len(mapTLDSubdomains))
	for z, sds := range mapTLDSubdomains {
		go func(zoneName string, subDomains []string) {
			defer wg.Done()

			logger.Logger.Debugf("Start update DNS record `%v` for zone name %v", strings.Join(subDomains, ","), zoneName)
			_ = cfdns.UpsertZoneDNSRecords(cfClient, ctx, zoneName, subDomains, ipv4, config.TTL)
		}(z, sds)
	}
	wg.Wait()
}
