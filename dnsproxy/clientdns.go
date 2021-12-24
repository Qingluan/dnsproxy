package main

import (
	"flag"

	"github.com/Qingluan/dnsproxy"
)

var (
	ServerMODE          = false
	ClientListenDNSPort = 0
	ListenAddr          = ""
)

func main() {
	flag.StringVar(&ListenAddr, "r", "0.0.0.0:60053", "\n\tserver : set  listen addr  \n\tclient: set remote's server addr")
	flag.IntVar(&ClientListenDNSPort, "p", 60053, "set client dns server's udp listen port ")
	flag.BoolVar(&ServerMODE, "s", false, "true to start as server dns proxy.")
	flag.Parse()

	cmdChan := make(chan string, 3)
	if ServerMODE {
		dnsproxy.NewDNSProxyServer(ClientListenDNSPort)
	} else {
		dnsproxy.NewDNSClientServer(ClientListenDNSPort, ListenAddr, cmdChan, nil)
	}
}
