package dnsproxy

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/miekg/dns"
)

var (
	SIG_CLEAN       = "[SIG_CLEAN]"
	SIG_EXIT        = "[SIG_EXIT]"
	DefaultLocalDNS = "223.5.5.5:53"
	l               = sync.RWMutex{}
)

// SetLocalDNS set default local dns : defualt is 223.5.5.5
func SetLocalDNS(dnsServer string) {
	l.Lock()
	defer l.Unlock()
	DefaultLocalDNS = dnsServer
}

func LocalQueryDNS(query *dns.Msg, localdnsserver ...string) (replyData []byte, err error) {
	c := new(dns.Client)

	if localdnsserver != nil {

		log.Println("[query "+localdnsserver[0]+"]:", query.Question[0].Name)
		reply, _, err := c.Exchange(query, localdnsserver[0])
		if err != nil {
			log.Println("[query local err]:", err)
			return nil, err
		}

		return reply.Pack()
	} else {
		log.Println("[query "+DefaultLocalDNS+"]:", query.Question[0].Name)
		reply, _, err := c.Exchange(query, DefaultLocalDNS)
		if err != nil {
			log.Println("[query local err]:", err)
			return nil, err
		}
		return reply.Pack()
	}
}

// Routine to handle inputs to Proxy port
func ClientProxy(listenPort int, cmdChan chan string, isLocalHost func(host string) bool, sendFund func(sendBuf []byte, otherDNSServer string) (reply []byte, err error)) (err error) {
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		log.Println("init udp port err:", err)
		return
	}
	udpListener, err := net.ListenUDP("udp", saddr)
	if err != nil {
		return err
	}
	for {
		otherDNS := ""
	label:
		select {
		case cmdMsg := <-cmdChan:
			if cmdMsg == SIG_EXIT {
				log.Println("exit this dns client")
				break label
			} else if cmdMsg == SIG_CLEAN {
				CleanCache()
			} else {
				otherDNS = cmdMsg
			}
		default:

			var buffer [1500]byte
			n, clientAddr, err := udpListener.ReadFromUDP(buffer[0:])
			if checkreport(1, err) {
				continue
			}
			nbuf := make([]byte, n)
			copy(nbuf, buffer[:n])

			go func(dnsData []byte, clientAddr *net.UDPAddr) {
				m := new(dns.Msg)
				if err := m.Unpack(nbuf); err == nil {
					if len(m.Question) > 0 {
						if reply, found := FindCache(m.Question[0].Name); found {
							replyMsg := new(dns.Msg)

							replyMsg.Unpack(reply.data)
							replyMsg.Id = m.Id
							if len(replyMsg.Answer) > 0 {
								log.Printf("local  (%5d)[%s] : %s (%d) \n", m.Id, m.Question[0].Name, replyMsg.Answer[0].String(), replyMsg.Id)
							} else {
								log.Printf("failed (%5d)[%s] : %s (%d) \n", m.Id, m.Question[0].Name, replyMsg.String(), replyMsg.Id)

							}
							data, _ := replyMsg.Pack()
							udpListener.WriteToUDP(data, clientAddr)
						} else {
							toLocal := false
							if isLocalHost != nil {
								toLocal = isLocalHost(m.Question[0].Name)
							}

							var replyData []byte
							var err error
							if toLocal {
								replyData, err = LocalQueryDNS(m)

							} else {
								replyData, err = sendFund(dnsData, otherDNS)
							}
							replyMsg := new(dns.Msg)
							replyMsg.Unpack(replyData)
							if len(replyMsg.Answer) > 0 {
								log.Printf("remote (%5d)[%s] : %s (%d) \n", m.Id, m.Question[0].Name, replyMsg.Answer[0].String(), replyMsg.Id)
							} else {
								log.Printf("failed (%5d)[%s] : %s (%d) \n", m.Id, m.Question[0].Name, replyMsg.String(), replyMsg.Id)
							}

							if err != nil {
								log.Println("dns remote resolve err:", err)
								return
							}
							_, err = udpListener.WriteToUDP(replyData, clientAddr)
							if err != nil {
								log.Println("reply dns err:", err)
								// return
							}
							if len(replyMsg.Answer) > 0 {
								RegistDNS(m.Question[0].Name, replyData)
							}

						}
					} else {
						log.Println("no query !")
					}

				} else {
					log.Println("not dns data jump")
				}
			}(nbuf, clientAddr)
		}

	}
}
