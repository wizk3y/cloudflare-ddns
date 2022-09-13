package domain_utils_test

import (
	"fmt"
	"testing"

	domain_utils "cloudflare-ddns/pkg/domain-utils"
)

func Test_ParseDomain(t *testing.T) {
	domains := []string{"consul.miraku.xyz", "miraku.xyz", "mirakustudio.com", "omr.mirakustudio.com", "www.test.librebee.com"}

	for _, d := range domains {
		sub, tld := domain_utils.ParseDomain(d)

		fmt.Println(sub, tld)
	}
}

func Test_GetMapTopLevelSubdomains(t *testing.T) {
	domains := []string{"consul.miraku.xyz", "miraku.xyz", "mirakustudio.com", "omr.mirakustudio.com", "www.librebee.com", "www.test.librebee.com", "consul.miraku.xyz"}

	mapTLDSubdomains := domain_utils.GetMapTopLevelSubdomains(domains)

	fmt.Println(mapTLDSubdomains)
}
