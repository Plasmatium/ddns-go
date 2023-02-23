package utils

import (
	"fmt"
	"testing"
)

func TestGetIP(t *testing.T) {
	ip, err := GetIP()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ip)
}
