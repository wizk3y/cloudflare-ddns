package domain_utils

import (
	"fmt"
	"testing"
)

func Test_GetTopLevelDomains(t *testing.T) {
	domains := []string{"consul.miraku.xyz", "miraku.xyz", "mirakustudio.com", "omr.mirakustudio.com"}

	topLevelDomains := GetTopLevelDomains(domains)

	fmt.Println(topLevelDomains)
}
