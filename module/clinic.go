package module

import (
	"errors"
	"xingren/model"
	"xingren/utils"
	"fmt"
)

type Clinic struct {
	ID         int64  `json:"id,omitempty"`
	DetailID   int64  `json:"detail_id,omitempty"`
	Date       string `json:"date,omitempty"`
	Morning    string `json:"morning,omitempty"`
	Afternoon  string `json:"afternoon,omitempty"`
	Night      string `json:"night,omitempty"`
	CreatedOn  int64  `json:"created_on,omitempty"`
	ModifiedOn int64  `json:"modified_on,omitempty"`
}

func ClinicAllByDetailIds(detailIds []int64) ([]Clinic, error) {
	list := make([]*model.Clinic, 0)
	where := []interface{}{[]interface{}{"queues_id", "in", detailIds}}
	if err := model.SelectData(&list, where, -1, 0, "", ""); err != nil {
		return nil, errors.New("get detail clinic list failed")
	}

	var clinicList []Clinic
	for _, l := range list {
		clinicList = append(clinicList, Clinic{
			ID:         l.ID,
			DetailID:   l.DetailID,
			Date:       l.Date,
			Morning:    l.Morning,
			Afternoon:  l.Afternoon,
			Night:      l.Night,
			CreatedOn:  l.CreatedOn,
			ModifiedOn: l.ModifiedOn,
		})
	}
	return clinicList, nil
}

func (c *Clinic) Detail() error {
	clinic := &model.Clinic{}
	if err := model.SelectOne(clinic, model.Clinic{ID: c.ID}); err != nil {
		return errors.New("get detail clinic failed")
	}
	if clinic.ID == 0 {
		return errors.New("detail clinic not exists")
	}

	c.ID = clinic.ID
	c.DetailID = clinic.DetailID
	c.Date = clinic.Date
	c.Morning = clinic.Morning
	c.Afternoon = clinic.Afternoon
	c.Night = clinic.Night
	c.CreatedOn = clinic.CreatedOn
	c.ModifiedOn = clinic.ModifiedOn
	return nil
}

func (c *Clinic) Total(keyword string, detailId int64) (int64, error) {
	total, err := model.Total(&model.Clinic{}, c.parseWhereConds(keyword, detailId))
	if err != nil {
		return 0, errors.New("get detail clinic count failed")
	}
	return total, nil
}

func (c *Clinic) Exists() (bool, error) {
	var where []interface{}
	if c.DetailID != 0 {
		where = append(where, []interface{}{"detail_id", "=", c.DetailID})
	}
	if c.Date != "" {
		where = append(where, []interface{}{"date", "=", c.Date})
	}
	if len(where) == 0 {
		return false, errors.New("no data found")
	}
	exit, err := model.Exit(&model.Clinic{}, where)
	if err != nil {
		return false, err
	}
	return exit, nil
}

func (c *Clinic) List(keyword string, detailId, page, pageSize int64) ([]Clinic, error) {
	list := make([]*model.Clinic, pageSize)
	if err := model.SelectData(&list, c.parseWhereConds(keyword, detailId), page, pageSize, "", ""); err != nil {
		return nil, errors.New("get queue robot list failed")
	}

	var clinic []Clinic
	for _, l := range list {
		clinic = append(clinic, Clinic{
			ID:         l.ID,
			DetailID:   c.DetailID,
			Date:       l.Date,
			Morning:    l.Morning,
			Afternoon:  l.Afternoon,
			Night:      l.Night,
			CreatedOn:  l.CreatedOn,
			ModifiedOn: l.ModifiedOn,
		})
	}
	return clinic, nil
}

func (c *Clinic) CreateOrUpdate() error {
	detail := &Detail{
		ID: c.DetailID,
	}
	if err := detail.Detail(); err != nil {
		return fmt.Errorf("detail error")
	}

	clinic := &model.Clinic{
		DetailID:  c.DetailID,
		Date:      c.Date,
		Morning:   c.Morning,
		Afternoon: c.Afternoon,
		Night:     c.Night,
	}
	updateField, err := utils.GormColumnToAllName(clinic)
	if err != nil {
		return err
	}
	where := map[string]interface{}{"date": c.Date, "detail_id": detail.ID}
	if c.ID != 0 {
		where = map[string]interface{}{"detail_id": detail.ID, "id": c.ID}
	}
	if err := model.UpSetData(clinic, where, updateField, true); err != nil {
		return errors.New("detail clinic update failed")
	}
	return nil
}

func (c *Clinic) Delete() error {
	if err := model.DelData(&model.Clinic{}, model.Clinic{ID: c.ID}); err != nil {
		return errors.New("detail clinic delete failed")
	}
	return nil
}

func (c *Clinic) parseWhereConds(keyword string, queuesId int64) []interface{} {
	var where []interface{}
	where = append(where, []interface{}{"detail_id", "=", queuesId})
	if keyword != "" {
		where = append(where, []interface{}{"date", "LIKE", fmt.Sprintf("%%%s%%", keyword)})
	}
	return where
}
