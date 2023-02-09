package main

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	_, err := GetIP()
	if err != nil {
		t.Fatal(err)
	}
}
