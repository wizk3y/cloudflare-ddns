package ip_utils

import (
	"fmt"
	"testing"
)

func Test_GetCurrentIPv6(t *testing.T) {
	fmt.Println(GetCurrentIPv6())
}
