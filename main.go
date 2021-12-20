package main

import (
	"fmt"
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
// func RunProxy(udpListener *net.UDPConn, sendFund func(sendBuf []byte) (reply []byte, err error)) {
func RunProxy(sendFund func(sendBuf []byte) (reply []byte, err error)) {

	var buffer [1500]byte
	for {
		n, _, err := ProxyConn.ReadFromUDP(buffer[0:])
		if checkreport(1, err) {
			continue
		}

		m := new(dns.Msg)
		if err := m.Unpack(buffer[:n]); err != nil {

		} else {
			fmt.Println("paresed:\n", m.String())
			fmt.Println("------------ end ------------")
		}
		// go func() {
		// 	if reply, err := sendFund(buffer[:n]); err != nil {
		// 		log.Println("[fail]:", err)
		// 	} else {
		// 		udpListener.WriteToUDP(reply, cliaddr)
		// 	}

		// }()
		// if !found {
		// 	conn = NewConnection(ServerAddr, cliaddr)
		// 	if conn == nil {
		// 		dunlock()
		// 		continue
		// 	}
		// 	ClientDict[saddr] = conn
		// 	dunlock()
		// 	Vlogf(2, "Created new connection for client %s\n", saddr)
		// 	// Fire up routine to manage new connection
		// 	go RunConnection(conn)
		// } else {
		// 	Vlogf(5, "Found connection for client %s\n", saddr)
		// 	dunlock()
		// }
		// Relay to server
		// _, err = conn.ServerConn.Write(buffer[0:n])
		// if checkreport(1, err) {
		// 	continue
		// }
	}
}

func main() {
	setup("baidu.com:80", 53)
	RunProxy(nil)
}
