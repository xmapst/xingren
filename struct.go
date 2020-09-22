package main

type Detail struct {
	DoctorID     int64      `json:"doctor_id,omitempty"`
	ID           int64      `json:"id,omitempty"`
	Title        string     `json:"title,omitempty"`
	Name         string     `json:"name,omitempty"`
	JobTitle     string     `json:"job_title,omitempty"`
	Department   string     `json:"department,omitempty"`
	Divide       string     `json:"divide,omitempty"`
	HospitalName string     `json:"hospital_name,omitempty"`
	AlmondNo     string     `json:"almond_no,omitempty"`
	Likes        string     `json:"likes,omitempty"`
	ViewNo       string     `json:"view_no,omitempty"`
	Introduction string     `json:"introduction,omitempty"`
	Clinic       []Clinic   `json:"clinic,omitempty"`
	Document     []Document `json:"document,omitempty"`
}

type Document struct {
	Title string `json:"title,omitempty"`
	Url   string `json:"url,omitempty"`
}

type Clinic struct {
	Date      string `json:"date,omitempty"`
	Morning   string `json:"morning,omitempty"`
	Afternoon string `json:"afternoon,omitempty"`
	Night     string `json:"night,omitempty"`
}

type TmpRequest struct {
	Html       string `json:"html,omitempty"`
	Success    bool   `json:"success,omitempty"`
	Code       int64  `json:"code"`
	Msg        string `json:"msg"`
	ErrCode    int64  `json:"errCode"`
	ErrMessage string `json:"errMessage"`
}