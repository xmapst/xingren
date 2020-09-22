package model

type Detail struct {
	Model
	ID           int64  `gorm:"column:id;primary_key" json:"id,omitempty"`
	DoctorID     int64  `gorm:"column:doctor_id;type:int(11);not null;default:0" json:"doctor_id,omitempty"`
	Title        string `gorm:"column:title;type:varchar(1000);not null;default:''" json:"title,omitempty"`
	Name         string `gorm:"column:name;type:varchar(1000);not null;default:''" json:"name,omitempty"`
	JobTitle     string `gorm:"column:job_title;type:varchar(1000);not null;default:''" json:"job_title,omitempty"`
	Department   string `gorm:"column:department;type:varchar(1000);not null;default:''" json:"department,omitempty"`
	Divide       string `gorm:"column:divide;type:varchar(1000);not null;default:''" json:"divide,omitempty"`
	HospitalName string `gorm:"column:hospital_name;type:varchar(1000);not null;default:''" json:"hospital_name,omitempty"`
	AlmondNo     string `gorm:"column:almond_no;type:varchar(1000);not null;default:''" json:"almond_no,omitempty"`
	Likes        string `gorm:"column:likes;type:varchar(1000);not null;default:''" json:"likes,omitempty"`
	ViewNo       string `gorm:"column:view_no;type:varchar(1000);not null;default:''" json:"view_no,omitempty"`
	Introduction string `gorm:"column:introduction;type:varchar(5000);not null;default:''" json:"introduction,omitempty"`
}

type Document struct {
	Model
	ID       int64  `gorm:"column:id;primary_key" json:"id,omitempty"`
	DetailID int64  `gorm:"column:detail_id;type:int(11);not null;default:0" json:"detail_id,omitempty"`
	Title    string `gorm:"column:title;type:varchar(5000);not null;default:''" json:"title,omitempty"`
	Url      string `gorm:"column:url;type:varchar(1000);not null;default:''" json:"url,omitempty"`
}

type Clinic struct {
	Model
	ID        int64  `gorm:"column:id;primary_key" json:"id,omitempty"`
	DetailID  int64  `gorm:"column:detail_id;type:int(11);not null;default:0" json:"detail_id,omitempty"`
	Date      string `gorm:"column:date;type:varchar(1000);not null;default:''" json:"date,omitempty"`
	Morning   string `gorm:"column:morning;type:varchar(1000);not null;default:''" json:"morning,omitempty"`
	Afternoon string `gorm:"column:afternoon;type:varchar(1000);not null;default:''" json:"afternoon,omitempty"`
	Night     string `gorm:"column:night;type:varchar(1000);not null;default:''" json:"night,omitempty"`
}
