package proxydns

import (
	"log"
	"net"

	"github.com/miekg/dns"
)

func ServerParseDNS(buffer []byte, replyFunc func(replyData []byte) error) (err error) {
	queryMsg := new(dns.Msg)
	err = queryMsg.Unpack(buffer)
	if err != nil {
		log.Println("not dns msg err:", err)
		return
	}
	c := new(dns.Client)
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	replyMsg, _, err := c.Exchange(queryMsg, net.JoinHostPort(config.Servers[0], config.Port))
	if err != nil {
		log.Println("resolve from "+net.JoinHostPort(config.Servers[0], config.Port), " err:", err)
		return
	}
	replyData, err := replyMsg.Pack()
	if err != nil {
		log.Println("[server] pack to data err :", err)
		return
	}
	err = replyFunc(replyData)
	return
}
