package crypt

import (
	"bytes"
)

func pkcs7Padding(data []byte, blockSize int) []byte {
	n := len(data)
	if n == 0 || blockSize < 1 {
		return data
	}
	paddingSize := blockSize - n%blockSize
	paddingText := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)

	return append(data, paddingText...)
}

func pkcs7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
