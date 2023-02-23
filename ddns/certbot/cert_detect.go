package certbot

import "time"

func DectectCertExpiration(domain, port string) (notBefore, notAfter time.Time, ok bool) {
	addr := domain + ":" + port
	println(addr)
	return
}