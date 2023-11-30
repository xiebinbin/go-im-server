package tt

import (
	"fmt"
	"imsdk/pkg/encrypt"
)

func Encrypt() {
	key := "1234567887654321"
	src := "helloworld"
	res, err := encrypt.AesGcmEncrypt([]byte(src), key)
	fmt.Println(res, err)
}
