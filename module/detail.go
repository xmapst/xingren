package module

import (
	"errors"
	"fmt"
	"xingren/model"
	"xingren/utils"
)

type Detail struct {
	ID           int64  `json:"id,omitempty"`
	DoctorID     int64  `json:"doctor_id,omitempty"`
	Title        string `json:"title,omitempty"`
	Name         string `json:"name,omitempty"`
	JobTitle     string `json:"job_title,omitempty"`
	Department   string `json:"department,omitempty"`
	Divide       string `json:"divide,omitempty"`
	HospitalName string `json:"hospital_name,omitempty"`
	AlmondNo     string `json:"almond_no,omitempty"`
	Likes        string `json:"likes,omitempty"`
	ViewNo       string `json:"view_no,omitempty"`
	Introduction string `json:"introduction,omitempty"`
	CreatedOn    int64  `json:"created_on,omitempty"`
	ModifiedOn   int64  `json:"modified_on,omitempty"`
}

func ListByIds(queuesIds []int64) ([]Detail, error) {
	s := &Detail{}
	return s.List(queuesIds, "", 0, 999)
}

func (d *Detail) List(detailIds []int64, keyword string, page, pageSize int64) ([]Detail, error) {
	list := make([]*model.Detail, pageSize)
	if err := model.SelectData(&list, d.parseWhereCond(detailIds, keyword), page, pageSize, "", ""); err != nil {
		return nil, errors.New("get detail list failed")
	}
	var detailSlice []Detail
	for _, detail := range list {
		detailSlice = append(detailSlice, Detail{
			ID:           detail.ID,
			DoctorID:     detail.DoctorID,
			Title:        detail.Title,
			Name:         detail.Name,
			JobTitle:     detail.JobTitle,
			Department:   detail.Department,
			Divide:       detail.Divide,
			HospitalName: detail.HospitalName,
			AlmondNo:     detail.AlmondNo,
			Likes:        detail.Likes,
			ViewNo:       detail.ViewNo,
			Introduction: detail.Introduction,
			CreatedOn:    detail.CreatedOn,
			ModifiedOn:   detail.ModifiedOn,
		})
	}
	return detailSlice, nil
}

func (d *Detail) Total(detailIds []int64, keyword string) (int64, error) {
	total, err := model.Total(&model.Detail{}, d.parseWhereCond(detailIds, keyword))
	if err != nil {
		return 0, fmt.Errorf("get detail list count failed, %q", err)
	}
	return total, nil
}

func (d *Detail) Exists() (bool, error) {
	var where []interface{}
	if d.Name != "" {
		where = append(where, []interface{}{"name", "=", d.Name})
	}
	if d.ID != 0 {
		where = append(where, []interface{}{"id", "!=", d.ID})
	}
	if d.DoctorID != 0 {
		where = append(where, []interface{}{"doctor_id", "!=", d.DoctorID})
	}
	if len(where) == 0 {
		return false, errors.New("no data found")
	}
	count, err := model.Total(&model.Detail{}, where)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *Detail) Delete() error {
	qwhere := model.Detail{ID: d.ID}
	if err := model.DelData(&model.Detail{}, qwhere); err != nil {
		return fmt.Errorf("delete detail failed %q", err)
	}
	dWhere := model.Document{DetailID: d.ID}
	if err := model.DelData(&model.Document{}, dWhere); err != nil {
		return fmt.Errorf("delete detail document failed %q", err)
	}
	cWhere := model.Clinic{DetailID: d.ID}
	if err := model.DelData(&model.Clinic{}, cWhere); err != nil {
		return fmt.Errorf("delete detail clinic failed %q", err)
	}
	return nil
}

func (d *Detail) Detail() error {
	detail := &model.Detail{}
	where := model.Detail{ID: d.ID}
	if err := model.SelectOne(detail, where); err != nil {
		return errors.New("get detail detail failed")
	}
	if detail.ID == 0 {
		return errors.New("queues detail not exists")
	}
	d.ID = detail.ID
	d.Name = detail.Name
	d.Title = detail.Title
	d.DoctorID = detail.DoctorID
	d.JobTitle = detail.JobTitle
	d.Department = detail.Department
	d.Divide = detail.Divide
	d.HospitalName = detail.HospitalName
	d.AlmondNo = detail.AlmondNo
	d.Likes = detail.Likes
	d.ViewNo = detail.ViewNo
	d.Introduction = detail.Introduction
	d.CreatedOn = detail.CreatedOn
	d.ModifiedOn = detail.ModifiedOn
	return nil
}

func (d *Detail) CreateOrUpdate() error {
	detail := &model.Detail{
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

	updateField, err := utils.GormColumnToAllName(detail)
	if err != nil {
		return err
	}

	where := map[string]interface{}{"id": d.ID}
	if err := model.UpSetData(detail, where, updateField, true); err != nil {
		return errors.New("detail update failed")
	}
	if d.ID == 0 {
		detail := &model.Detail{}
		where := model.Detail{ID: d.ID}
		if err := model.SelectOne(detail, where); err != nil {
			return errors.New("queues update failed")
		}
		if detail.ID == 0 {
			return errors.New("queues update failed")
		}
		d.ID = detail.ID
	}
	return nil
}

func (d *Detail) SearchName() error {
	detail := &model.Detail{}
	where := model.Detail{Name: d.Name, DoctorID: d.DoctorID}
	if err := model.SelectOne(detail, where); err != nil {
		return errors.New("detail does not exist")
	}
	if detail.ID == 0 {
		return errors.New("detail does not exist")
	}
	d.ID = detail.ID
	return nil
}

func (d *Detail) SearchDoctorId() error {
	detail := &model.Detail{}
	where := model.Detail{DoctorID: d.DoctorID}
	if err := model.SelectOne(detail, where); err != nil {
		return errors.New("detail does not exist")
	}
	if detail.ID == 0 {
		return errors.New("detail does not exist")
	}
	d.ID = detail.ID
	return nil
}

func (d *Detail) parseWhereCond(queuesIds []int64, keyword string) []interface{} {
	var where []interface{}
	if keyword != "" {
		where = append(where, []interface{}{"name", "LIKE", fmt.Sprintf("%%%s%%", keyword)})
	}
	where = append(where, []interface{}{"id", "IN", queuesIds})
	return where
}
