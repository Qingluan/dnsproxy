package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

type Connection struct {
	ClientAddr *net.UDPAddr // Address of the client
	ServerConn *net.UDPConn // UDP connection to server
}

var ProxyConn *net.UDPConn

// Address of server
var ServerAddr *net.UDPAddr

// Mapping from client addresses (as host:port) to connection
var ClientDict map[string]*Connection = make(map[string]*Connection)

// Generate a new connection by opening a UDP connection to the server
func NewConnection(srvAddr, cliAddr *net.UDPAddr) *Connection {
	conn := new(Connection)
	conn.ClientAddr = cliAddr
	srvudp, err := net.DialUDP("udp", nil, srvAddr)
	if checkreport(1, err) {
		return nil
	}
	conn.ServerConn = srvudp
	return conn
}

func setup(hostport string, port int) bool {
	// Set up Proxy
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if checkreport(1, err) {
		return false
	}
	pudp, err := net.ListenUDP("udp", saddr)
	if checkreport(1, err) {
		return false
	}
	ProxyConn = pudp
	Vlogf(2, "Proxy serving on port %d\n", port)

	// Get server address
	srvaddr, err := net.ResolveUDPAddr("udp", hostport)
	if checkreport(1, err) {
		return false
	}
	ServerAddr = srvaddr
	Vlogf(2, "Connected to server at %s\n", hostport)
	return true
}

// Routine to handle inputs to Proxy port
func ClientProxy(listenPort int, sendFund func(sendBuf []byte) (reply []byte, err error)) (err error) {
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", listenPort))
	udpListener, err := net.ListenUDP("udp", saddr)
	if err != nil {
		return err
	}
	var buffer [1500]byte
	for {
		n, cliaddr, err := ProxyConn.ReadFromUDP(buffer[0:])
		if checkreport(1, err) {
			continue
		}

		m := new(dns.Msg)
		if err := m.Unpack(buffer[:n]); err == nil {

			go func() {
				if reply, err := sendFund(buffer[:n]); err != nil {
					log.Println("[fail]:", err)
				} else {
					udpListener.WriteToUDP(reply, cliaddr)
				}

			}()
		} else {
			log.Println("not dns data jump")
		}
	}
}

func main() {
	port := 6053
	// setup("baidu.com:80", 53)
	ClientProxy(port, func(sendBuf []byte) (reply []byte, err error) {
		return
	})
}
