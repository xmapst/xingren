package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"xingren/model"
	"xingren/parser"
	"xingren/persist"
	"xingren/proxypool"
	//"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"log"
	"net/http"
)

func init() {
	model.Setup()
	model.InitDatabase()
}

const (
	detailUrl   = "http://xingren.com/wx/doctor/%d/feeds"
	documentUrl = "http://xingren.com/page/home/%d/documents"
	clinicUrl   = "http://xingren.com/page/home/%d/clinic"
)

// 重试策略是爬到数据为止
func main() {
	// 启动数据保存的协程
	go persist.ItemServer()
	// 初始化detail控制器
	detailCollector := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(1),
		// 使用异步并发
		colly.Async(true),
		// 打印debug日志
		//colly.Debugger(&debug.LogDebugger{}),
	)

	// 使用随机UA请求头
	extensions.RandomUserAgent(detailCollector)

	// 使用返回头的Referer
	extensions.Referer(detailCollector)

	// 禁用Keep Alives, 忽略https证书验证
	detailCollector.WithTransport(&http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	// 控制并发数量
	err := detailCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 300})
	if err != nil {
		log.Panic(err)
	}
	// 使用代理
	rp, err := proxy.RoundRobinProxySwitcher(proxypool.IpProt...)
	if err != nil {
		log.Fatal(err)
	}
	detailCollector.SetProxyFunc(rp)

	// 异常重试
	detailCollector.OnError(func(r *colly.Response, err error) {
		// 出现错误时重试
		if err != nil {
			err = r.Request.Retry()
			if err != nil {
				log.Println(err)
			}
		}
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		log.Println("visiting:", r.URL)
	})

	// 克隆document 及 clinic 的控制器
	documentCollector := detailCollector.Clone()
	clinicCollector := detailCollector.Clone()

	// 处理详情页数据
	detailCollector.OnResponse(func(r *colly.Response) {
		log.Println("Got:", r.Request.URL.String())
		// 返回非正常数据重新爬取
		tmpRequest := TmpRequest{}
		if err := json.Unmarshal(r.Body, &tmpRequest); err == nil && !tmpRequest.Success {
			err = r.Request.Retry()
			if err != nil {
				log.Println(r.Request.URL, err)
			}
			return
		}
		// 处理正常返回的数据
		id, err := parser.ParseDetail(r.Body, r.Request.URL.String(), persist.DetailChannel)
		if err != nil {
			log.Println(r.Request.URL, err)
			return
		}
		// 爬取对应的文章页面
		err = documentCollector.Visit(fmt.Sprintf(documentUrl, id))
		if err != nil {
			log.Println(err)
		}
		// 爬取对应的坐诊页面
		err = clinicCollector.Visit(fmt.Sprintf(clinicUrl, id))
		if err != nil {
			log.Println(err)
		}
	})
	// 处理文章页数据
	documentCollector.OnResponse(func(r *colly.Response) {
		log.Println("Got:", r.Request.URL.String())
		// 返回非正常数据重新爬取
		tmpRequest := TmpRequest{}
		if err := json.Unmarshal(r.Body, &tmpRequest); err == nil && !tmpRequest.Success {
			err = r.Request.Retry()
			if err != nil {
				log.Println(r.Request.URL, err)
			}
			return
		}
		// 处理正常返回的数据
		err := parser.ParseDocument(r.Body, r.Request.URL.String(), persist.DocumentChannel)
		if err != nil {
			log.Println(r.Request.URL, err)
			return
		}
	})

	// 处理坐诊页数据
	clinicCollector.OnResponse(func(r *colly.Response) {
		log.Println("Got:", r.Request.URL.String())
		// 返回非正常数据重新爬取
		tmpRequest := TmpRequest{}
		if err := json.Unmarshal(r.Body, &tmpRequest); err == nil && !tmpRequest.Success {
			err = r.Request.Retry()
			if err != nil {
				log.Println(r.Request.URL, err)
			}
			return
		}
		// 处理正常返回的数据
		err := parser.ParseClinic(r.Body, r.Request.URL.String(), persist.ClinicChannel)
		if err != nil {
			log.Println(r.Request.URL, err)
			return
		}
	})

	// 添加爬虫队列
	for i := 1; i <= 500000; i++ {
		err := detailCollector.Visit(fmt.Sprintf(detailUrl, i))
		if err != nil {
			log.Println(err)
		}
	}
	// 等待爬取完成
	detailCollector.Wait()
	documentCollector.Wait()
	clinicCollector.Wait()

	// 读入之前爬取的内容， 已弃用
	//var wg sync.WaitGroup
	//f, err := os.Open("xingren.json")
	//if err != nil {
	//	log.Panic(err)
	//}
	//defer f.Close()
	//reader := bufio.NewReader(f)
	//for {
	//	line, err := reader.ReadString('\n')
	//	if err != nil || io.EOF == err {
	//		break
	//	}
	//	t := &Detail{}
	//	if err := json.Unmarshal([]byte(line), &t); err != nil {
	//		fmt.Println(line)
	//		log.Panic(err)
	//		//continue
	//	}
	//	wg.Add(1)
	//	func(t *Detail) {
	//		defer wg.Done()
	//		detail := &module.Detail{
	//			ID:           t.ID,
	//			DoctorID:     t.DoctorID,
	//			Title:        t.Title,
	//			Name:         t.Name,
	//			JobTitle:     t.JobTitle,
	//			Department:   t.Department,
	//			Divide:       t.Divide,
	//			HospitalName: t.HospitalName,
	//			AlmondNo:     t.AlmondNo,
	//			Likes:        t.Likes,
	//			ViewNo:       t.ViewNo,
	//			Introduction: t.Introduction,
	//		}
	//		if err := detail.CreateOrUpdate(); err != nil {
	//			log.Println(err)
	//		}
	//		if t.Document != nil {
	//			for _, d := range t.Document {
	//				document := &module.Document{
	//					DetailID: t.ID,
	//					Title:    d.Title,
	//					Url:      d.Url,
	//				}
	//				if err := document.CreateOrUpdate(); err != nil {
	//					log.Println(err)
	//				}
	//			}
	//		}
	//
	//		if t.Clinic != nil {
	//			for _, c := range t.Clinic {
	//				clinic := &module.Clinic{
	//					DetailID:  t.ID,
	//					Date:      c.Date,
	//					Morning:   c.Morning,
	//					Afternoon: c.Afternoon,
	//					Night:     c.Night,
	//				}
	//				if err := clinic.CreateOrUpdate(); err != nil {
	//					log.Println(err)
	//				}
	//			}
	//		}
	//	}(t)
	//}
	//wg.Wait()
}
