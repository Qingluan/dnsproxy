package proxydns

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

// Routine to handle inputs to Proxy port
func ClientProxy(listenPort int, sendFund func(sendBuf []byte) (reply []byte, err error)) (err error) {
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		log.Println("init udp port err:", err)
		return
	}
	udpListener, err := net.ListenUDP("udp", saddr)
	if err != nil {
		return err
	}
	var buffer [1500]byte
	for {
		n, clientAddr, err := udpListener.ReadFromUDP(buffer[0:])
		if checkreport(1, err) {
			continue
		}

		m := new(dns.Msg)
		if err := m.Unpack(buffer[:n]); err == nil {
			if len(m.Question) > 0 {
				log.Printf("query (%d) : %s \n", len(m.Question), m.Question[0].Name)
				// m.Question[0].Name
				if reply, found := FindCache(m.Question[0].Name); found {

					udpListener.WriteToUDP(reply.data, clientAddr)
				} else {
					go func(host string, senddata []byte, clientAddr *net.UDPAddr) {
						replyData, err := sendFund(senddata)
						if err != nil {
							log.Println("dns remote resolve err:", err)
							return
						}

						_, err = udpListener.WriteToUDP(replyData, clientAddr)
						if err != nil {
							log.Println("reply dns err:", err)
							// return
						}
						RegistDNS(host, replyData)
					}(m.Question[0].Name, buffer[:n], clientAddr)
				}
			} else {
				log.Println("no query !")
			}

		} else {
			log.Println("not dns data jump")
		}
	}
}
