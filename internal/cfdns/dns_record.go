package cfdns

import (
	"cloudflare-ddns/internal/logger"
	"context"
	"fmt"
	"sync"

	"github.com/cloudflare/cloudflare-go"
)

func UpsertZoneDNSRecords(cfClient *cloudflare.API, ctx context.Context, zoneName string, subdomains []string, ipv4 string, ttl int) error {
	// validate ttl
	if ttl != 1 && ttl < 60 {
		ttl = 60
	}

	if ttl > 86400 {
		ttl = 86400
	}

	zoneID, mapNameRecord, err := getZoneAndRecords(cfClient, ctx, zoneName)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	wg.Add(len(subdomains))
	for _, sd := range subdomains {
		go func(name string) {
			defer wg.Done()

			r := fmt.Sprintf("%s.%s", name, zoneName)
			if name == "@" {
				r = zoneName
			}

			// check record existed
			oldRecord, ok := mapNameRecord[r]
			if !ok || oldRecord.ZoneName != zoneName {
				// create new record
				logger.Logger.Debugf("Start create DNS record for name %v of zone %v", name, zoneName)
				_ = insertDNSRecordIP(cfClient, ctx, zoneID, name, ipv4, ttl)
				return
			}

			// check content of old record is identical with current
			if oldRecord.Content == ipv4 {
				logger.Logger.Infof("Record name %v of zone %v already point to ip %v", name, zoneName, ipv4)
				return
			}

			logger.Logger.Debugf("Start update DNS record for name %v of zone %v", name, zoneName)
			_ = updateDNSRecordIP(cfClient, ctx, zoneID, oldRecord.ID, ipv4, ttl)
		}(sd)
	}
	wg.Wait()

	return nil
}

func getZoneAndRecords(cfClient *cloudflare.API, ctx context.Context, zoneName string) (string, map[string]cloudflare.DNSRecord, error) {
	zoneID, err := cfClient.ZoneIDByName(zoneName)
	if err != nil {
		logger.Logger.Errorf("Error when get zone id by name. Zone name %v, error details: %v", zoneName, err)
		return "", nil, err
	}

	var dnsRecords []cloudflare.DNSRecord
	dnsRecords, err = cfClient.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{
		Type: "A",
	})
	if err != nil {
		logger.Logger.Errorf("Error when get dns record of zone. Zone name %v, error details: %v", zoneName, err)
		return zoneID, nil, err
	}

	mapNameRecord := make(map[string]cloudflare.DNSRecord)
	for _, r := range dnsRecords {
		mapNameRecord[r.Name] = r
	}

	return zoneID, mapNameRecord, nil
}

func insertDNSRecordIP(cfClient *cloudflare.API, ctx context.Context, zoneID, name, newIP string, ttl int) error {
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
		TTL:     ttl,
	})
	if err != nil {
		logger.Logger.Errorf("Error when add dns record. Name %v, zone ID %v, error details: %v", name, zoneID, err)
		return err
	}

	logger.Logger.Infof("New record created. Zone ID %v - name %v pointed to %v", zoneID, name, newIP)

	return nil
}

func updateDNSRecordIP(cfClient *cloudflare.API, ctx context.Context, zoneID, recordID, newIP string, ttl int) error {
	if recordID == "" {
		return nil
	}

	if newIP == "" {
		return nil
	}

	err := cfClient.UpdateDNSRecord(ctx, zoneID, recordID, cloudflare.DNSRecord{
		Content: newIP,
		TTL:     ttl,
	})
	if err != nil {
		logger.Logger.Errorf("Error when update dns record. Record ID %v, zone ID %v, error details: %v", recordID, zoneID, err)
		return err
	}

	logger.Logger.Infof("Record ID: %v updated. Zone ID %v pointed to %v", recordID, zoneID, newIP)

	return nil
}
