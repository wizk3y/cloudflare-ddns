package domainutil

import (
	"fmt"
	"strings"
)

func ParseDomain(d string) (subdomain, topLevelDomain string) {
	domainLevels := strings.Split(d, ".")

	if len(domainLevels) < 2 {
		// invalid domain
		return "", ""
	}

	return strings.Join(domainLevels[:len(domainLevels)-2], "."), strings.Join(domainLevels[len(domainLevels)-2:], ".")
}

func GetMapTopLevelSubdomains(domains []string) map[string][]string {
	var mapTLDSubdomainsUniq = make(map[string]map[string]bool)

	for _, d := range domains {
		s, tld := ParseDomain(d)

		if len(tld) == 0 {
			continue
		}

		if _, ok := mapTLDSubdomainsUniq[tld]; !ok {
			mapTLDSubdomainsUniq[tld] = make(map[string]bool)
		}

		if len(s) == 0 {
			s = "@"
		}

		mapTLDSubdomainsUniq[tld][s] = true
	}

	var mapTLDSubdomains = make(map[string][]string)
	for tld, mapUniq := range mapTLDSubdomainsUniq {
		for s := range mapUniq {
			mapTLDSubdomains[tld] = append(mapTLDSubdomains[tld], s)
		}
	}

	return mapTLDSubdomains
}

func BuildListDomain(tld string, subDomains []string) []string {
	listDomains := make([]string, 0)
	for _, sd := range subDomains {
		if sd == "@" {
			listDomains = append(listDomains, tld)
		} else {
			listDomains = append(listDomains, fmt.Sprintf("%s.%s", sd, tld))
		}
	}

	return listDomains
}
