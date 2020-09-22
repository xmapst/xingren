package proxypool

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestProxy(t *testing.T) {
	var wg sync.WaitGroup
	var ipAvailable []string
	for _, ip := range IpProt {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			var speed, status = ProxyTest(ip)
			if status == 200 {
				log.Println(speed, status, ip)
				ipAvailable = append(ipAvailable, ip)
			}
		}(ip)

	}
	wg.Wait()
	fmt.Println(ipAvailable)
	t.Log("pass")
}
