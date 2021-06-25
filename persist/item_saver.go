package persist

import (
	"log"
	"xingren/model"
	"xingren/module"
)

var (
	DetailChannel   = make(chan model.Detail, 65535)
	ClinicChannel   = make(chan []model.Clinic, 65535)
	DocumentChannel = make(chan []model.Document, 65535)
)

func ItemServer() {
	for {
		select {
		case d := <-DetailChannel:
			detail := &module.Detail{
				ID:           d.ID,
				DoctorID:     d.DoctorID,
				Title:        d.Title,
				Name:         d.Name,
				JobTitle:     d.JobTitle,
				Department:   d.Department,
				Divide:       d.Divide,
				HospitalName: d.HospitalName,
				AlmondNo:     d.AlmondNo,
				Likes:        d.Likes,
				ViewNo:       d.ViewNo,
				Introduction: d.Introduction,
			}
			if err := detail.CreateOrUpdate(); err != nil {
				log.Println(err)
			}
		case c := <-ClinicChannel:
			for _, c := range c {
				clinic := &module.Clinic{
					DetailID:  c.DetailID,
					Date:      c.Date,
					Morning:   c.Morning,
					Afternoon: c.Afternoon,
					Night:     c.Night,
				}
				if err := clinic.CreateOrUpdate(); err != nil {
					log.Println(err)
				}
			}
		case d := <-DocumentChannel:
			for _, d := range d {
				document := &module.Document{
					DetailID: d.DetailID,
					Title:    d.Title,
					Url:      d.Url,
				}
				if err := document.CreateOrUpdate(); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
