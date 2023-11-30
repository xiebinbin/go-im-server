package applet

import (
	"encoding/base64"
	"imsdk/pkg/encrypt"
	"imsdk/pkg/funcs"
	"strconv"
)

func getCode() string {
	t := funcs.GetMillis()
	token := "72FDA07287D854DF7EBD56AC2BA2E7DDMnRnbWQzYWlvZnlr"
	str := token + "|" + strconv.Itoa(int(t))
	key := "ghizg3wwxtg5bhfl"
	uid := "855ad8681d0d"
	encStr := encrypt.AesCbcEncrypt([]byte(str), key)
	str1 := uid + "|" + "3" + "|" + encStr
	return base64.StdEncoding.EncodeToString([]byte(str1))
}
