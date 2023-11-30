package common

import (
	"imsdk/internal/common/dao/common/region"
	"imsdk/pkg/redis"
	"strconv"
)

type CountryListResp struct {
	ID       string `bson:"_id" json:"id"`
	Name     string `bson:"name" json:"name"`
	Children int8   `json:"children"`
}

const (
	CityListCacheTag = "chat:im:cityList:ver"
)

func GetCountryList(lang string) ([]CountryListResp, error) {
	var res []CountryListResp
	data, _ := region.New().GetCountryList(lang)
	if len(data) == 0 {
		return res, nil
	}
	tempData := make(map[string]string, 0)
	for _, v := range data {
		tempData[v.PId] = v.ID
	}
	for _, v := range data {
		children := 0
		if _, ok := tempData[v.ID]; ok {
			children = 1
		}
		if v.Level == 1 {
			resTmp := CountryListResp{
				ID:       v.ID,
				Name:     v.Name,
				Children: int8(children),
			}
			res = append(res, resTmp)
		}
	}
	return res, nil
}

func GetCityList(lang string, version int) ([]region.Region, int, error) {
	result, err := redis.Client.Get(CityListCacheTag).Result()
	if err != nil && err != redis.NilErr {
		return []region.Region{}, 0, err
	}
	ver, _ := strconv.Atoi(result)
	if ver == 0 {
		ver = 1
	}
	if version >= ver {
		return []region.Region{}, 0, nil
	}
	data, _ := region.New().GetListByLang(lang, "_id,code,lang,level,name,pid")
	redis.Client.Set(CityListCacheTag, ver, 0)
	return data, ver, nil
}
