package proxypool

import (
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var IpProt = []string{
	"http://221.182.31.54:8080",
	"http://221.182.31.54:8080",
	"http://211.137.52.159:8080",
	"http://80.241.222.138:80",
	"http://120.79.186.104:8118",
	"http://122.226.57.70:8888",
	"http://113.100.209.215:3128",
	"http://113.214.13.1:1080",
	"http://106.52.153.210:3128",
	"http://47.113.121.239:3128",
	"http://221.13.156.158:55443",
	"http://27.36.118.2:3128",
	"http://47.113.121.239:3128",
	"http://113.214.13.1:1080",
	"http://115.223.7.110:80",
	"http://211.137.52.159:8080",
	"http://39.106.155.65:8888",
	"http://39.106.223.134:80",
	"http://47.110.163.118:80",
	"http://60.255.151.82:80",
	"http://211.137.52.158:8080",
	"http://47.107.108.83:3128",
	"http://39.106.155.65:8888",
	"http://39.106.155.65:8888",
	"http://116.196.85.150:3128",
	"http://116.196.85.150:3128",
	"http://116.196.85.150:3128",
	"http://47.100.111.190:9999",
	"http://60.255.151.82:80",
	"http://221.180.170.104:8080",
	"http://210.26.49.89:3128",
	"http://210.26.49.89:3128",
	"http://161.202.226.194:80",
	"http://210.26.49.89:3128",
	"http://113.214.13.1:1080",
	"http://218.60.8.99:3129",
	"http://111.229.43.213:8118",
	"http://47.113.121.239:3128",
	"http://59.36.10.79:3128",
	"http://14.139.62.246:80",
	"http://122.224.65.197:3128",
	"http://80.241.222.137:80",
	"http://116.231.96.146:8118",
	"http://123.56.161.63:80",
	"http://222.66.94.130:80",
	"http://39.106.223.134:80",
	"http://113.214.13.1:1080",
	"http://47.113.121.239:3128",
	"http://218.60.8.83:3129",
	"http://60.255.151.82:80",
	"http://218.60.8.99:3129",
	"http://116.196.85.150:3128",
	"http://113.214.13.1:1080",
	"http://223.82.106.253:3128",
	"http://91.205.174.26:80",
	"http://221.180.170.104:8080",
	"http://218.75.102.198:8000",
	"http://106.3.45.16:58080",
	"http://69.197.181.202:3128",
	"http://222.66.94.130:80",
	"http://36.103.223.96:3128",
	"http://211.137.52.158:8080",
	"http://47.103.110.169:80",
	"http://58.250.21.56:3128",
	"http://47.103.110.169:80",
	"http://173.192.128.238:8123",
	"http://169.57.157.146:8123",
	"http://125.73.220.18:49128",
	"http://159.8.114.37:80",
	"http://222.85.28.130:52590",
	"http://122.224.65.197:3128",
	"http://220.194.226.136:3128",
	"http://47.94.200.124:3128",
	"http://59.37.18.243:3128",
	"http://221.6.201.18:9999",
	"http://61.178.149.237:59042",
	"http://60.214.134.54:45947",
	"http://222.175.171.6:8080",
	"http://49.234.135.40:8085",
	"http://222.175.171.6:8080",
}

func Random() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return IpProt[r.Intn(len(IpProt))]
}

func ProxyTest(proxy_addr string) (Speed int, Status int) {
	//检测代理iP访问地址
	var testUrl string
	//判断传来的代理IP是否是https
	if strings.Contains(proxy_addr, "https") {
		testUrl = "https://icanhazip.com"
	} else {
		testUrl = "http://icanhazip.com"
	}
	// 解析代理地址
	proxy, err := url.Parse(proxy_addr)
	//设置网络传输
	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	// 创建连接客户端
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	begin := time.Now() //判断代理访问时间
	// 使用代理IP访问测试地址
	res, err := httpClient.Get(testUrl)

	if err != nil {
		//log.Println(err)
		return
	}
	defer res.Body.Close()
	speed := int(time.Now().Sub(begin).Nanoseconds() / 1000 / 1000) //ms
	//判断是否成功访问，如果成功访问StatusCode应该为200
	if res.StatusCode != http.StatusOK {
		//log.Println(err)
		return
	}
	return speed, res.StatusCode
}
