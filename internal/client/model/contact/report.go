package contact

import (
	contact2 "imsdk/internal/common/dao/contact"
)

type ReportContact struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
type ReportResp struct {
	UID   string `json:"uid"`
	Phone string `json:"phone"`
}
type ReportRequest struct {
	Phones []ReportContact `json:"phones"`
}

type PhoneInfo struct {
	Prefix   string `json:"prefix"`
	Phone    string `json:"phone"`
	OldPhone string `json:"old_phone"`
	Alias    string `json:"alias"`
}

func ExistsMineInfo(uid, prefix, phone string) contact2.Contact {
	data := contact2.New().GetIsMineInfo(uid, prefix, phone)
	return data
}
