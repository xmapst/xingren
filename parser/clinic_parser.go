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

// 获取门诊时间
// http://xingren.com/page/home/%d/clinic
func ParseClinic(contents []byte, u string, out chan []model.Clinic) (err error) {
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
	request := clinicFilter(dom, detailId)
	if len(request) == 0 {
		return errors.New("data is empty")
	}
	out <- request
	return nil
}

func clinicFilter(dom *goquery.Document, detailId int64) (result []model.Clinic) {
	dom.Find("table > tbody > tr").Each(func(i int, trNodes *goquery.Selection) {
		if trNodes == nil {
			return
		}
		tr := strconv.Itoa(i + 1)
		trHtml, _ := trNodes.Find("tr:nth-child(" + tr + ") > th").Html()
		strS := utils.ChangeSpan(trHtml)
		result = append(result, model.Clinic{
			DetailID:  detailId,
			Date:      utils.DelSpace(strS[0]),
			Morning:   utils.DelSpace(trNodes.Find("tr:nth-child(" + tr + ") > td:nth-child(2)").Text()),
			Afternoon: utils.DelSpace(trNodes.Find("tr:nth-child(" + tr + ") > td:nth-child(3)").Text()),
			Night:     utils.DelSpace(trNodes.Find("tr:nth-child(" + tr + ") > td:nth-child(4)").Text()),
		})
	})
	return result
}
