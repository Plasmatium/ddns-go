package main

import (
	"fmt"
	"os"
	"sync"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	log "github.com/sirupsen/logrus"
)

var client *alidns20150109.Client
var initClientOnce sync.Once
var domainName string

func init() {
	client = mustCreateClient()
	domainName = os.Getenv("DOMAIN_NAME")
	log.WithField("domain_name", domainName).Info("domain name loaded")
}

func mustCreateClient() *alidns20150109.Client {
	initClientOnce.Do(func() {
		config := &openapi.Config{
			// Protocol:        Ref("http"),
			AccessKeyId:     Ref(os.Getenv("ACCESS_KEY_ID")),
			AccessKeySecret: Ref(os.Getenv("ACCESS_KEY_SECRET")),
		}
		config.Endpoint = Ref("alidns.cn-shanghai.aliyuncs.com")
		var err error
		if client, err = alidns20150109.NewClient(config); err != nil {
			log.Fatal(err)
		}
	})
	return client
}

func getDNSRecord(recordKeyword string) (recordID, prevIP string, ok bool) {
	req := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName:  &domainName,
		RRKeyWord:   Ref(recordKeyword),
		TypeKeyWord: Ref("A"),
	}
	result, err := client.DescribeDomainRecordsWithOptions(req, &util.RuntimeOptions{})
	if err != nil {
		log.WithError(err).Error("get recordID failed")
		return
	}

	records := result.Body.DomainRecords.Record
	if len(records) == 0 {
		err = fmt.Errorf("records not found")
		log.WithError(err).
			WithField("keyword", recordKeyword).
			Error("get recordID failed")
		return
	}
	record := records[0]

	return *record.RecordId, *record.Value, true
}

func setDNS(recordID, recordKeyword, ip string) {
	req := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: &recordID,
		RR:       &recordKeyword,
		Type:     Ref("A"),
		Value:    &ip,
	}
	runtime := &util.RuntimeOptions{}

	_, err := client.UpdateDomainRecordWithOptions(req, runtime)
	if err != nil {
		log.WithError(err).Error("failed to set dns")
	} else {
		log.WithField("record", recordKeyword).
			WithField("ip", ip).
			Info("set dns record success")
	}
}

func TrySetDNS(recordKeyword string) {
	recordID, prevIP, ok := getDNSRecord(recordKeyword)
	if !ok {
		log.WithField("record_name", recordKeyword).Info("setting on new record")
	}

	ip, err := GetIP()
	if err != nil {
		log.WithError(err).Error("try set dns failed")
		return
	}

	if prevIP == ip {
		log.WithField("ip", ip).
			Debug("ip needs to set is same as prev ip, no operation needs to be done")
		return
	}

	setDNS(recordID, recordKeyword, ip)
}
