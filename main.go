package main

import (
	"cloudflare-ddns/pkg/config"
	domain_utils "cloudflare-ddns/pkg/domain-utils"
	ip_utils "cloudflare-ddns/pkg/ip-utils"
	"cloudflare-ddns/pkg/log"
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"github.com/robfig/cron"
)

var (
	// this like an state when running in cron mode
	currentIP string
)

func init() {
	log.InitLogger()
}

func main() {
	if config.Mode == "single" {
		pointToCurrentIP()
	} else {
		c := cron.New()
		err := c.AddFunc("0 */5 * * * *", pointToCurrentIP)
		if err != nil {
			log.Logger.Fatalf("Error when add cron function, details: %v", err)
		}

		log.Logger.Infof("Cron function added successfully.")

		c.Start()
		select {}
	}
}

func pointToCurrentIP() {
	var ipv4 = ip_utils.GetCurrentIPv4()
	// var ipv6 = ip_utils.GetCurrentIPv6()

	log.Logger.Infof("Current IP v4: %v", ipv4)

	// check is IP changed
	if ipv4 == currentIP {
		log.Logger.Debugf("The current IP of the machine has not changed.")
		return
	}

	// Construct a new API object
	cfClient, err := cloudflare.New(config.CfApiKey, config.CfApiEmail)
	if err != nil {
		log.Logger.Errorf("Error when create cf client, details: %v", err)
		return
	}

	// Most API calls require a Context
	ctx := context.Background()

	// aggregate domains list
	mapTLDSubdomains := domain_utils.GetMapTopLevelSubdomains(config.Domains)

	// get list zone
	for z, sds := range mapTLDSubdomains {
		log.Logger.Debugf("Start update DNS record for zone name: %v", z)
		_ = upsertZoneDNSRecords(cfClient, ctx, z, sds, ipv4)
	}
}

func upsertZoneDNSRecords(cfClient *cloudflare.API, ctx context.Context, zoneName string, subdomains []string, ipv4 string) error {
	zoneID, mapNameRecord, err := getZoneAndRecords(cfClient, ctx, zoneName)
	if err != nil {
		return err
	}

	for _, sd := range subdomains {
		r := fmt.Sprintf("%s.%s", sd, zoneName)
		if sd == "@" {
			r = zoneName
		}

		// check record existed
		oldRecord, ok := mapNameRecord[r]
		if !ok || oldRecord.ZoneName != zoneName {
			// create new record
			log.Logger.Debugf("Start create DNS record for name: %v", sd)
			_ = insertDNSRecordIP(cfClient, ctx, zoneID, sd, ipv4)
			continue
		}

		// check content of old record is identical with current
		if oldRecord.Content == ipv4 {
			log.Logger.Infof("Record name %v already point to ip %v", sd, ipv4)
			continue
		}

		log.Logger.Debugf("Start update DNS record for name: %v", sd)
		_ = updateDNSRecordIP(cfClient, ctx, zoneID, oldRecord.ID, ipv4)
	}

	return nil
}

func getZoneAndRecords(cfClient *cloudflare.API, ctx context.Context, zoneName string) (string, map[string]cloudflare.DNSRecord, error) {
	zoneID, err := cfClient.ZoneIDByName(zoneName)
	if err != nil {
		log.Logger.Errorf("Error when get zone id by name, details: %v", err)
		return "", nil, err
	}

	var dnsRecords []cloudflare.DNSRecord
	dnsRecords, err = cfClient.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{
		Type: "A",
	})
	if err != nil {
		log.Logger.Errorf("Error when get dns record of zone, details: %v", err)
		return zoneID, nil, err
	}

	mapNameRecord := make(map[string]cloudflare.DNSRecord)
	for _, r := range dnsRecords {
		mapNameRecord[r.Name] = r
	}

	return zoneID, mapNameRecord, nil
}

func insertDNSRecordIP(cfClient *cloudflare.API, ctx context.Context, zoneID, name, newIP string) error {
	if name == "" {
		return nil
	}

	if newIP == "" {
		return nil
	}

	_, err := cfClient.CreateDNSRecord(ctx, zoneID, cloudflare.DNSRecord{
		Type:    "A",
		Name:    name,
		Content: newIP,
		TTL:     1,
	})
	if err != nil {
		log.Logger.Errorf("Error when add dns record, details: %v", err)
		return err
	}

	log.Logger.Infof("New record created. Zone ID %v - name %v pointed to %v", zoneID, name, newIP)

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
