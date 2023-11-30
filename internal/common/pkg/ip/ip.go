package ip

import (
	"encoding/json"
	"imsdk/pkg/errno"
	"imsdk/pkg/redis"
	"imsdk/pkg/request"
	"time"
)

const parseIpUrl = "http://ip-api.com/json/"

type Info struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	Timezone    string `json:"timezone"`
}

func GetIpInfo(ip string) (Info, error) {
	cacheTag := "tmm:im:ip:" + ip
	cacheRes, err := redis.Client.Get(cacheTag).Result()
	if err != nil && err != redis.NilErr {
		return Info{}, errno.Add("sys err", errno.SysErr)
	}
	var ipInfo = Info{}
	if cacheRes != "" {
		err1 := json.Unmarshal([]byte(cacheRes), &ipInfo)
		if err1 != nil {
			return Info{}, nil
		}
	} else {
		url := parseIpUrl + ip
		response, err1 := request.Get(url, map[string]string{"SocketData-Type": "Application/json"})
		if err1 != nil {
			return Info{}, nil
		}
		err1 = json.Unmarshal(response, &ipInfo)
		if err1 != nil {
			return Info{}, nil
		}
		redis.Client.Set(cacheTag, string(response), time.Second*180)
	}
	return ipInfo, nil
}
