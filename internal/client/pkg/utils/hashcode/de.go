package hashcode

import "github.com/speps/go-hashids/v2"

func De(group string, code string) uint {
	id := uint(0)
	hd := hashids.NewData()
	hd.Salt = group
	hd.MinLength = 4
	h, err := hashids.NewWithData(hd)
	if err == nil {
		tmp := h.Decode(code)
		id = uint(tmp[0])
	}
	return id
}
