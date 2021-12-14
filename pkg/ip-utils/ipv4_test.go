package ip_utils

import (
	"fmt"
	"testing"
)

func Test_GetCurrentIPv4(t *testing.T) {
	fmt.Println(GetCurrentIPv4())
}
