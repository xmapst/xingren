package parser

import (
	"bytes"
	"errors"
	"xingren/model"
	"xingren/utils"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

// 获取文章
// http://xingren.com/page/home/%d/documents
func ParseDocument(contents []byte, u string, out chan []model.Document) (err error) {
	id := regexp.MustCompile(idRe).FindAllString(u, -1)
	if len(id) == 0 {
		return errors.New("id is empty")
	}
	detailId, err := strconv.ParseInt(id[0], 10, 64)
	if err != nil {
		return err
	}
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
	if err != nil {
		return err
	}
	request := documentFilter(dom, detailId)
	if len(request) == 0 {
		return errors.New("data is empty")
	}
	out <- request
	return nil
}

func documentFilter(dom *goquery.Document, detailId int64) (result []model.Document) {
	dom.Find("ul.doc-list > li").Each(func(_ int, liNode *goquery.Selection) {
		if liNode == nil {
			return
		}
		a, _ := liNode.Find("a").Attr("href")
		title := utils.DelSpace(liNode.Text())
		result = append(result, model.Document{
			DetailID: detailId,
			Title:    utils.DelSpace(title),
			Url:      a,
		})
	})
	return result
}
