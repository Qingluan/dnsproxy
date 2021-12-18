package main

import (
	"sync"
	"time"
)

var (
	lock    = sync.RWMutex{}
	Cachsed = make(map[string]Cache)
)

type Cache struct {
	data   []byte
	create time.Time
}

func RegistDNS(host string, replyDNS []byte) {
	if replyDNS == nil || len(replyDNS) == 0 {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	c := Cache{
		data:   replyDNS,
		create: time.Now(),
	}
	Cachsed[host] = c
}

func FindCache(host string) (c Cache, found bool) {
	c, found = Cachsed[host]
	return
}
