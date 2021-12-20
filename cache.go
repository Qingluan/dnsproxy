package dnsproxy

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
	// client *net.UDPAddr
}

func CleanCache() {
	lock.Lock()
	defer lock.Unlock()
	Cachsed = make(map[string]Cache)
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
		// client: client,
	}
	Cachsed[host] = c
}

func FindCache(host string) (c Cache, found bool) {
	c, found = Cachsed[host]
	return
}
