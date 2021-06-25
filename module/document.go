package module

import (
	"errors"
	"fmt"
	"xingren/model"
	"xingren/utils"
)

type Document struct {
	ID         int64  `json:"id,omitempty"`
	DetailID   int64  `json:"detail_id,omitempty"`
	Title      string `json:"title,omitempty"`
	Url        string `json:"url,omitempty"`
	CreatedOn  int64  `json:"created_on,omitempty"`
	ModifiedOn int64  `json:"modified_on,omitempty"`
}

func DocumentAllByDetailIds(detailIds []int64) ([]Document, error) {
	list := make([]*model.Document, 0)
	where := []interface{}{[]interface{}{"detail_id", "in", detailIds}}
	if err := model.SelectData(&list, where, -1, 0, "", ""); err != nil {
		return nil, errors.New("get detail document list failed")
	}

	var documentList []Document
	for _, l := range list {
		documentList = append(documentList, Document{
			ID:         l.ID,
			DetailID:   l.DetailID,
			Title:      l.Title,
			Url:        l.Url,
			CreatedOn:  l.CreatedOn,
			ModifiedOn: l.ModifiedOn,
		})
	}
	return documentList, nil
}

func (d *Document) Detail() error {
	document := &model.Document{}
	if err := model.SelectOne(document, model.Document{ID: d.ID}); err != nil {
		return errors.New("get queue robot detail failed")
	}
	if document.ID == 0 {
		return errors.New("queues robot detail not exists")
	}

	d.ID = document.ID

	return nil
}

func (d *Document) Total(keyword string, detailId int64) (int64, error) {
	total, err := model.Total(&model.Document{}, d.parseWhereConds(keyword, detailId))
	if err != nil {
		return 0, errors.New("get detail document count failed")
	}
	return total, nil
}

func (d *Document) Exists() (bool, error) {
	var where []interface{}
	if d.DetailID != 0 {
		where = append(where, []interface{}{"detail_id", "=", d.DetailID})
	}
	if d.Title != "" {
		where = append(where, []interface{}{"name", "=", d.Title})
	}
	if len(where) == 0 {
		return false, errors.New("no data found")
	}
	exit, err := model.Exit(&model.Document{}, where)
	if err != nil {
		return false, err
	}
	return exit, nil
}

func (d *Document) List(keyword string, detailId, page, pageSize int64) ([]Document, error) {
	list := make([]*model.Document, pageSize)
	if err := model.SelectData(&list, d.parseWhereConds(keyword, detailId), page, pageSize, "", ""); err != nil {
		return nil, errors.New("get detail document list failed")
	}

	var document []Document
	for _, l := range list {
		document = append(document, Document{
			ID:         l.ID,
			Title:      l.Title,
			Url:        l.Url,
			CreatedOn:  l.CreatedOn,
			ModifiedOn: l.ModifiedOn,
		})
	}
	return document, nil
}

func (d *Document) CreateOrUpdate() error {
	detail := &Detail{
		ID: d.DetailID,
	}
	if err := detail.Detail(); err != nil {
		return fmt.Errorf("detail error")
	}

	document := &model.Document{
		DetailID: detail.ID,
		Title:    d.Title,
		Url:      d.Url,
	}
	updateField, err := utils.GormColumnToAllName(document)
	if err != nil {
		return err
	}
	where := map[string]interface{}{"title": d.Title, "detail_id": detail.ID}
	if d.ID != 0 {
		where = map[string]interface{}{"detail_id": detail.ID, "id": d.ID}
	}
	if err := model.UpSetData(document, where, updateField, true); err != nil {
		return errors.New("detail document update failed")
	}
	return nil
}

func (d *Document) Delete() error {
	if err := model.DelData(&model.Document{}, model.Document{ID: d.ID}); err != nil {
		return errors.New("detail document delete failed")
	}
	return nil
}

func (d *Document) parseWhereConds(keyword string, queuesId int64) []interface{} {
	var where []interface{}
	where = append(where, []interface{}{"detail_id", "=", queuesId})
	if keyword != "" {
		where = append(where, []interface{}{"name", "LIKE", fmt.Sprintf("%%%s%%", keyword)})
	}
	return where
}
