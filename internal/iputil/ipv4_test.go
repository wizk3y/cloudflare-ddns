package iputil_test

import (
	"fmt"
	"testing"

	"cloudflare-ddns/internal/iputil"
)

func Test_GetCurrentIPv4(t *testing.T) {
	fmt.Println(iputil.GetCurrentIPv4())
}
