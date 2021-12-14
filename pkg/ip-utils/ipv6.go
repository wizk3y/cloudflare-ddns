package ip_utils

import (
	"regexp"
	"strconv"
	"strings"
)

var ipv6Regexp = regexp.MustCompile(`(?m)^    inet6 ([0-9a-f:]+)`)

func GetCurrentIPv6() string {
	// res, err := exec.Command("ip", "-6", "addr", "list", "scope", "global", "-deprecated").Output()
	// if err != nil {
	// 	// not support
	// 	return ""
	// }

	res := []byte(`
en0: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
    options=400<CHANNEL_IO>
    ether 8c:85:90:7e:51:9c
    inet6 fe80::140c:9c32:89ca:113%en0 prefixlen 64 secured scopeid 0x4
    inet6 2402:800:61b1:e81a:1020:f3f7:92de:4a89 prefixlen 64 autoconf secured
    inet6 2402:800:61b1:e81a:5534:536f:d75b:8705 prefixlen 64 autoconf temporary
    inet6 2402:800:61b1:e81a::a prefixlen 64 dynamic
    inet 192.168.1.12 netmask 0xffffff00 broadcast 192.168.1.255
    nd6 options=201<PERFORMNUD,DAD>
    media: autoselect
    status: active
utun3: flags=8051<UP,POINTOPOINT,RUNNING,MULTICAST> mtu 1380
    inet6 fe80::37f0:209e:39f3:88d4%utun3 prefixlen 64 scopeid 0xe
    nd6 options=201<PERFORMNUD,DAD>
utun4: flags=8051<UP,POINTOPOINT,RUNNING,MULTICAST> mtu 1380
    inet6 fe80::ed44:57c1:aa4d:d999%utun4 prefixlen 64 scopeid 0xf
    nd6 options=201<PERFORMNUD,DAD>
utun5: flags=8051<UP,POINTOPOINT,RUNNING,MULTICAST> mtu 1380
    inet6 fe80::3c2c:ad17:78c6:26bc%utun5 prefixlen 64 scopeid 0x10
    nd6 options=201<PERFORMNUD,DAD>
	`)

	matches := ipv6Regexp.FindAllStringSubmatch(string(res), -1)

	for _, m := range matches {
		if len(m) < 2 {
			continue
		}

		if v := validate(m[1]); v {
			return m[1]
		}
	}

	return ""
}

func validate(ip string) bool {
	prefixValue, err := strconv.ParseInt(ip[:strings.Index(ip, ":")], 16, 32)
	if err != nil {
		return false
	}

	return (prefixValue & 0xfe00) == 0xfc00
}
