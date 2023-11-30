package common

import (
	"imsdk/internal/common/dao/common/countryphoneprefix"
	"imsdk/pkg/redis"
	"strconv"
	"strings"
)

type CountryPhonePrefixListResp struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Lang  string `json:"lang"`
	Code  string `json:"code"`
	CCode string `json:"c_code"`
}

const (
	CountryPhonePrefixLisCacheTag = "chat:im:phonePrefix:ver"
	CountryPhonePrefix            = "Chat:im:phonePrefix"
)

func CountryPhonePrefixList(lang string, version int) ([]CountryPhonePrefixListResp, int, error) {
	var res []CountryPhonePrefixListResp
	result, err := redis.Client.Get(CountryPhonePrefixLisCacheTag).Result()
	if err != nil && err != redis.NilErr {
		return res, 0, err
	}
	ver, _ := strconv.Atoi(result)
	if ver == 0 {
		ver = 1
	}
	if version >= ver {
		return []CountryPhonePrefixListResp{}, 0, nil
	}
	data, _ := countryphoneprefix.New().GetList(lang)
	if len(data) == 0 {
		data, _ = countryphoneprefix.New().GetList(GetDefaultLang())
	}
	for _, v := range data {
		resTmp := CountryPhonePrefixListResp{
			ID:    v.ID,
			Name:  v.Name,
			Code:  v.Code,
			Lang:  v.Lang,
			CCode: strings.ToUpper(v.CCode),
		}
		res = append(res, resTmp)
	}
	redis.Client.Set(CountryPhonePrefixLisCacheTag, ver, 0)
	return res, ver, nil
}

func PrefixIsExist(prefix string) bool {
	rs, err := redis.Client.SCard(CountryPhonePrefix).Result()
	if err != nil && err != redis.NilErr {
		return true
	}
	if rs == 0 {
		data, _ := countryphoneprefix.New().GetList("tr")
		for _, v := range data {
			redis.Client.SAdd(CountryPhonePrefix, v.Code)
		}
	}
	res, er := redis.Client.SIsMember(CountryPhonePrefix, prefix).Result()
	if er != nil && er != redis.NilErr {
		return true
	}
	return res
}
