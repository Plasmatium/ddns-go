package dns

import (
	"fmt"
	"testing"
)

func TestGetDNSRecords(*testing.T) {
	domainName = "yucy-top.love"
	fmt.Println(GetDNSRecord(domainName, "@"))
}