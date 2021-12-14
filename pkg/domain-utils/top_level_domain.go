package domain_utils

import "strings"

func GetTopLevelDomains(domains []string) []string {
	var mapTopLevelDomains = make(map[string]bool)

	for _, d := range domains {
		domainLevels := strings.Split(d, ".")

		if len(domainLevels) < 2 {
			// domain invalid
			continue
		}

		topLevelDomain := strings.Join(domainLevels[len(domainLevels)-2:], ".")

		mapTopLevelDomains[topLevelDomain] = true
	}

	topLevelDomains := make([]string, len(mapTopLevelDomains))

	for d := range mapTopLevelDomains {
		topLevelDomains = append(topLevelDomains, d)
	}

	return topLevelDomains
}
