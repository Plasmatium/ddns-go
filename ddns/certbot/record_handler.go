package certbot

import (
	"ddns-go/dns"
)

func AddChallengeRecord(domainName, txtValue string) {
	dns.AddDNSRecord(domainName, "", "_acme-challenge", "TXT", txtValue)
}

func DeleteChallengeRecord(domainName string) {
	records, _ := dns.GetDNSRecords(domainName, "_acme-challenge", "TXT")
	var ridList []string
	for _, r := range records {
		ridList = append(ridList, *r.RecordId)
	}
	dns.DeleteDNSRecords(ridList)
}

