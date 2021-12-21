package main

import (
	"flag"

	"github.com/Qingluan/dnsproxy"
)

func main() {
	p := 60053
	s := "127.0.0.1:60053"
	flag.IntVar(&p, "p", 60053, "set port ")
	flag.StringVar(&s, "s", "127.0.0.1:60053", "remote dns server proxy ")

	flag.Parse()
	cmdChan := make(chan string, 3)
	dnsproxy.NewDNSClientServer(p, s, cmdChan, nil)

}
