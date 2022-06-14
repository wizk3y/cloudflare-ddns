package main

import (
	ip_utils "cloudflare-ddns/pkg/ip-utils"
	"cloudflare-ddns/pkg/log"
	"context"
	"flag"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/robfig/cron"
)

var (
	cfApiKey   string
	cfApiEmail string
	domains    []string

	currentIP string
)

func init() {
	var d string

	flag.StringVar(&cfApiKey, "cf-api-key", "", "")
	flag.StringVar(&cfApiEmail, "cf-api-email", "", "")
	flag.StringVar(&d, "domains", "", "")
	flag.Parse()

	domains = strings.Split(d, ",")

	log.InitLogger()
}

func main() {
	c := cron.New()
	err := c.AddFunc("0 */5 * * * *", updateDomainDNS)
	if err != nil {
		log.Logger.Fatalf("Error when add cron function, details: %v", err)
	}

	log.Logger.Infof("Cron function added successfully.")

	c.Start()
	select {}
}

func updateDomainDNS() {
	var ipv4 = ip_utils.GetCurrentIPv4()
	// var ipv6 = ip_utils.GetCurrentIPv6()

	log.Logger.Infof("Current IP v4: %v", ipv4)

	// check is IP changed
	if ipv4 == currentIP {
		log.Logger.Debugf("The current IP of the machine has not changed.")
		return
	}

	// Construct a new API object
	cfClient, err := cloudflare.New(cfApiKey, cfApiEmail)
	if err != nil {
		log.Logger.Errorf("Error when create cf client, details: %v", err)
		return
	}

	// Most API calls require a Context
	ctx := context.Background()

	// get list zone
	for _, d := range domains {
		log.Logger.Debugf("Start update DNS record domain: %v", d)
		_ = updateDNSRecordDomain(cfClient, ctx, d, ipv4)
	}
}

func updateDNSRecordDomain(cfClient *cloudflare.API, ctx context.Context, zoneName, ipv4 string) error {
	zoneID, err := cfClient.ZoneIDByName(zoneName)
	if err != nil {
		log.Logger.Errorf("Error when get zone id by name, details: %v", err)
		return err
	}

	var dnsRecords []cloudflare.DNSRecord
	dnsRecords, err = cfClient.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{})
	if err != nil {
		log.Logger.Errorf("Error when get dns record of zone, details: %v", err)
		return err
	}

	for _, r := range dnsRecords {
		if r.Content == ipv4 {
			log.Logger.Infof("Zone %v already point to ip %v", zoneName, ipv4)
			continue
		}

		switch r.Type {
		case "A":
			_ = updateDNSRecordIP(cfClient, ctx, zoneID, r.ID, ipv4)
			// updateDNSRecordIP(cfClient, ctx, zoneID, r.ID, ipv6)
		case "AAAA":
			_ = updateDNSRecordIP(cfClient, ctx, zoneID, r.ID, ipv4)
			// updateDNSRecordIP(cfClient, ctx, zoneID, r.ID, ipv6)
		}
	}

	return nil
}

func updateDNSRecordIP(cfClient *cloudflare.API, ctx context.Context, zoneID, recordID, newIP string) error {
	if recordID == "" {
		return nil
	}

	if newIP == "" {
		return nil
	}

	err := cfClient.UpdateDNSRecord(ctx, zoneID, recordID, cloudflare.DNSRecord{
		Content: newIP,
	})
	if err != nil {
		log.Logger.Errorf("Error when update dns record, details: %v", err)
		return err
	}

	log.Logger.Infof("Record ID: %v updated. Zone ID %v pointed to %v", recordID, zoneID, newIP)

	return nil
}
