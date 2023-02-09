package main

import (
	"fmt"
	"os"
	"testing"
)

func TestGetDNSRecords(*testing.T) {
	os.Setenv("ACCESS_KEY_ID", "LTAI5tN9VwQq3mpfD6AeCaDS")
	os.Setenv("ACCESS_KEY_SECRET", "oJcSaqVo2CYWasJJCYiwiiGnqV8dfY")
	domainName = "yucy-top.love"
	// client = mustCreateClient()
	fmt.Println(getDNSRecord("@"))
}

