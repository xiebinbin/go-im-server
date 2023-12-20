package contact

import (
	"imsdk/internal/common/dao/user/contact"
	user2 "imsdk/internal/common/model/user"
	"imsdk/pkg/funcs"
)

type UpdateBaseRequest struct {
	UId        string `json:"uid" binding:"required"`
	Alias      string `json:"alias"`
	RemarkText string `json:"remark_text"`

	Phones    []map[string]interface{} `json:"phones"`
	Relations []map[string]interface{} `json:"relations"`
	Emails    []map[string]interface{} `json:"emails"`
	Dates     []map[string]interface{} `json:"dates"`
	Companies []map[string]interface{} `json:"companies"`
	Schools   []map[string]interface{} `json:"schools"`
	Address   []map[string]interface{} `json:"addrs"`
}

func GetContactsDetail(uid string, ids []string) []contact.Contact {
	data, _ := contact.New().GetContactsInfo(uid, ids)
	return data
}

func UpdateAlias(uid, objUId, alias string) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"alias": alias,
	}
	matchedCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if matchedCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}

func GetAlias(uid, objUId string) string {
	return contact.New().GetAlias(uid, objUId)
}

func UpdatePhone(uid, objUId, phone string) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"phone": phone,
	}
	modCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if modCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}

func addRemark(uid, objUId string, data map[string]interface{}) (string, error) {
	id := contact.GetId(uid, objUId)
	t := funcs.GetMillis()
	data["_id"] = id
	data["uid"] = uid
	data["obj_uid"] = objUId
	data["itime"] = t
	data["utime"] = t
	lastId, addErr := contact.New().AddRemark(data)
	return lastId, addErr
}

func UpdatePhones(uid, objUId string, phones []map[string]interface{}) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"phones": phones,
	}
	modCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if modCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}

func UpdateRelations(uid, objUId string, relations []map[string]interface{}) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"relations": relations,
	}
	modCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if modCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}

func UpdateEmails(uid, objUId string, emails []map[string]interface{}) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"emails": emails,
	}
	modCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if modCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}

func UpdateDates(uid, objUId string, dates []map[string]interface{}) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"dates": dates,
	}
	modCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if modCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}
func UpdateCompanies(uid, objUId string, companies []map[string]interface{}) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"companies": companies,
	}
	modCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if modCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}

func UpdateSchools(uid, objUId string, schools []map[string]interface{}) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"schools": schools,
	}
	modCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if modCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}

func UpdateAddress(uid, objUId string, Address []map[string]interface{}) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"addrs": Address,
	}
	modCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if modCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}

func UpdateRemarkText(uid, objUId, remark string) error {
	if err := user2.GetUserErr(objUId); err != nil {
		return err
	}
	upData := map[string]interface{}{
		"remark": remark,
	}
	modCount, err := contact.New().UpdateRemark(uid, objUId, upData)
	if modCount <= 0 {
		_, addErr := addRemark(uid, objUId, upData)
		return addErr
	}
	return err
}
