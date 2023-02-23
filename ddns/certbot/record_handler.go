package certbot

import (
	"ddns-go/dns"
	"ddns-go/utils"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

func addDomainRecord(domainName, txtValue string) {
	req := &alidns20150109.AddDomainRecordRequest{
		DomainName: &domainName,
		RR: utils.Ref("_acme-challenge"),
		Value: &txtValue,
		Type: utils.Ref("TXT"),
	}
	runtime := &util.RuntimeOptions{}
	if _, err := dns.GetSDKClient().AddDomainRecordWithOptions(req, runtime); err != nil {
		log.WithError(err).Error("failed to a")
	}
}

func deleteDomainRecord(domainName string) {

}