package hashcode

import "github.com/speps/go-hashids/v2"

func En(group string, id uint) string {
	var code string
	hd := hashids.NewData()
	hd.Salt = group
	hd.MinLength = 4
	h, err := hashids.NewWithData(hd)
	if err == nil {
		code, _ = h.Encode([]int{int(id)})
	}
	return code
}
