package dns

import (
	"fmt"
	"os"
	"sync"

	"ddns-go/utils"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	log "github.com/sirupsen/logrus"
)

var client *alidns20150109.Client
var initClientOnce sync.Once
var domainName string

func init() {
	domainName = os.Getenv("DOMAIN_NAME")
	log.WithField("domain_name", domainName).Info("domain name loaded")
}

func GetSDKClient() *alidns20150109.Client {
	initClientOnce.Do(func() {
		config := &openapi.Config{
			// Protocol:        Ref("http"),
			AccessKeyId:     utils.Ref(os.Getenv("ACCESS_KEY_ID")),
			AccessKeySecret: utils.Ref(os.Getenv("ACCESS_KEY_SECRET")),
		}
		config.Endpoint = utils.Ref("alidns.cn-shanghai.aliyuncs.com")
		var err error
		if client, err = alidns20150109.NewClient(config); err != nil {
			log.Fatal(err)
		}
	})
	return client
}

func GetDNSRecord(domainName, recordKeyword, rType string) (recordID, prevIP string, ok bool) {
	req := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName:  &domainName,
		RRKeyWord:   utils.Ref(recordKeyword),
		TypeKeyWord: &rType,
	}
	result, err := GetSDKClient().DescribeDomainRecordsWithOptions(req, &util.RuntimeOptions{})
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

func SetDNSRecord(domainName, recordID, recordKeyword, rType, value string) {
	req := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: &recordID,
		RR:       &recordKeyword,
		Type:     &rType,
		Value:    &value,
	}
	runtime := &util.RuntimeOptions{}

	if _, err := GetSDKClient().UpdateDomainRecordWithOptions(req, runtime); err != nil {
		log.WithError(err).Error("failed to set dns")
	} else {
		log.WithField("record", recordKeyword).
			WithField("ip", value).
			Info("set dns record success")
	}
}
func AddDNSRecord(domainName, recordID, recordKeyword, rType, value string) {
	req := &alidns20150109.AddDomainRecordRequest{
		DomainName: &domainName,
		RR:         &recordKeyword,
		Type:       &rType,
		Value:      &value,
	}
	runtime := &util.RuntimeOptions{}

	if _, err := GetSDKClient().AddDomainRecordWithOptions(req, runtime); err != nil {
		log.WithError(err).Error("failed to add dns")
	} else {
		log.WithField("record", recordKeyword).
			WithField("ip", value).
			Info("add dns record success")
	}
}

// TrySetDNS
// step 1. find self public ip
// step 2. get previous ip, if not exist, change updateDNS from setDNS to addDNS
// step 3. compare previous and target ip, if same, no need to update
// step 4. do update
func TrySetDNS(recordKeyword string) {
	ip, err := utils.GetIP()
	if err != nil {
		log.WithError(err).Error("try set dns failed on get self public ip")
		return
	}

	updateDNS := SetDNSRecord
	recordID, prevIP, ok := GetDNSRecord(domainName, recordKeyword, "A")
	if !ok {
		log.WithField("record_name", recordKeyword).Info("previous record not found, setting on new record")
		updateDNS = AddDNSRecord
	}

	if prevIP == ip {
		log.WithField("ip", ip).
			Debug("ip needs to set is same as prev ip, no operation needs to be done")
		return
	}

	updateDNS(domainName, recordID, recordKeyword, "A", ip)
}
