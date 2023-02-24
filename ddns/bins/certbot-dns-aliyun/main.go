package main

import (
	"ddns-go/certbot"
	"ddns-go/dns"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	// 操作类型，authenticator: 域名认证，添加 DNS TXT 记录; cleanup: 认证通过后，删除此 DNS TXT 记录
	op := flag.String("o", "auth", "Operate: authenticator or cleanup")
	confFile := flag.String("c", "credential.json", "Config file path")
	
	help := flag.Bool("h", false, "show help infomation")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	domain := os.Getenv("CERTBOT_DOMAIN")
	validation := os.Getenv("CERTBOT_VALIDATION")
	if domain == "" {
		log.Warn("domain is empty, please ensure called by certbot, nothing to do")
		return
	}
	log := log.WithField("domain", domain)

	loadConfig(*confFile)

	switch *op {
	case "auth":
		log.Info("adding txt record")
		certbot.AddChallengeRecord(domain, validation)
		time.Sleep(time.Second * 30)

	case "cleanup":
		log.Info("doing cleanup")
		certbot.DeleteChallengeRecord(domain)
	}
}

type authInfo struct {
	Ak string `json:"ACCESS_KEY_ID"`
	Sk string `json:"ACCESS_KEY_SECRET"`
}

func loadConfig(confFilePath string) {
	log := log.WithField("path", confFilePath)
	log.Info("try loading conf file")
	bs, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		log.WithError(err).WithField("path", confFilePath).Panic("failed to load config file")
	}

	var auth authInfo
	json.Unmarshal(bs, &auth)
	dns.MustRebuildClient(auth.Ak, auth.Sk)
}