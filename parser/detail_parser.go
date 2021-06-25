package parser

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
	"xingren/model"
	"xingren/utils"
)

// 解析详情页
// http://xingren.com/wx/doctor/%d/feeds
func ParseDetail(contents []byte, u string, out chan model.Detail) (id int64, err error) {
	doctorId := regexp.MustCompile(idRe).FindAllString(u, -1)
	if len(doctorId) == 0 {
		return 0, errors.New("id is empty")
	}
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
	if err != nil {
		return 0, err
	}
	request := detailFilter(dom)
	if request.Name == "" || request.ID == 0 {
		return 0, errors.New("doctor don't exist")
	}
	request.DoctorID, err = strconv.ParseInt(doctorId[0], 10, 64)
	if err != nil {
		return 0, err
	}
	out <- request
	return request.ID, nil
}

func detailFilter(dom *goquery.Document) (detail model.Detail) {
	detail.Title = utils.DelSpace(dom.Find("head > title").Text())

	// find <div class="profile-infos flexbox">
	dom.Find("header.profile > div.profile-infos > ul.flex").Each(func(_ int, liNodes *goquery.Selection) {
		liNodes.Find("li:nth-child(1)").Each(func(_ int, selection *goquery.Selection) {
			detail.Name = utils.DelSpace(selection.Find("span.name").Text())
			detail.Divide = utils.DelSpace(selection.Find("span.divide").Text())
		})
		liHtml, _ := liNodes.Find("li:nth-child(1)").Html()
		strS := utils.ChangeSpan(liHtml)
		if len(strS) != 0 && len(strS) == 5 {
			detail.JobTitle = utils.DelSpace(strS[2])
			detail.Department = utils.DelSpace(strS[4])
		}
		detail.HospitalName = utils.DelSpace(liNodes.Find("li:nth-child(2)").Text())
		detail.AlmondNo = utils.DelSpace(liNodes.Find("li:nth-child(3)").Text())
		detail.ID, _ = strconv.ParseInt(strings.Split(liNodes.Find("li:nth-child(3)").Text(), ":")[1], 10, 64)
	})

	dom.Find("header.profile > ul.profile-stats").Each(func(_ int, liNodes *goquery.Selection) {
		detail.Likes = liNodes.Find("li:nth-child(1)").Text()
		detail.ViewNo = liNodes.Find("li:nth-child(2)").Text()
	})

	detail.Introduction = utils.DelSpace(utils.DelSpace(dom.Find("div.panel-body").Text()))
	return detail
}
