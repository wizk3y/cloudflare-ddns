package domainutil_test

import (
	"fmt"
	"testing"

	"cloudflare-ddns/internal/domainutil"
)

func Test_ParseDomain(t *testing.T) {
	domains := []string{"consul.miraku.xyz", "miraku.xyz", "mirakustudio.com", "omr.mirakustudio.com", "www.test.librebee.com"}

	for _, d := range domains {
		sub, tld := domainutil.ParseDomain(d)

		fmt.Println(sub, tld)
	}
}

func Test_GetMapTopLevelSubdomains(t *testing.T) {
	domains := []string{"consul.miraku.xyz", "miraku.xyz", "mirakustudio.com", "omr.mirakustudio.com", "www.librebee.com", "www.test.librebee.com", "consul.miraku.xyz"}

	mapTLDSubdomains := domainutil.GetMapTopLevelSubdomains(domains)

	fmt.Println(mapTLDSubdomains)
}
