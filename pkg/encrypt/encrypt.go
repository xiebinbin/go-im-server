package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"imsdk/pkg/funcs"
	"math"
	"strings"
)

func Padding(src []byte, blockSize int) []byte {
	// calculate padding length
	n := blockSize - len(src)%blockSize
	// Fill n of the original clear text with n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	return append(src, temp...)
}

func UnPadding(cipherText []byte) []byte {
	end := cipherText[len(cipherText)-1]
	return cipherText[:len(cipherText)-int(end)]
}

func AesCbcEncrypt(src []byte, key string) string {
	block, _ := aes.NewCipher([]byte(key))
	src = Padding(src, block.BlockSize())
	iv := []byte(funcs.GetRandString(16))
	BlockMode := cipher.NewCBCEncrypter(block, iv)
	resText := make([]byte, len(src))
	BlockMode.CryptBlocks(resText, src)
	resText = funcs.BytesCombine(iv, resText)
	res := base64.URLEncoding.EncodeToString(resText)
	return strings.TrimRight(string(res), "=")
}

func AesCbcDecrypt(src, key string) []byte {
	// base64 url decode
	length := len(src)
	if length%4 != 0 {
		padNUm := int(math.Ceil(float64(length)/4)*4) - length
		for i := 0; i < padNUm; i++ {
			src = src + "="
		}
	}
	decodeBytes, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return nil
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil
	}
	iv := decodeBytes[:16]
	dataBytes := decodeBytes[16:]
	blockMode := cipher.NewCBCDecrypter(block, iv)
	oriData := make([]byte, len(dataBytes))
	blockMode.CryptBlocks(oriData, dataBytes)
	return UnPadding(oriData)
}

func AesGcmEncrypt(src []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	iv := []byte(funcs.GetRandString(12))
	//tag := []byte(funcs.GetRandString(16))
	dst := make([]byte, 0)
	tx := aesgcm.Seal(nil, iv, src, nil)
	fmt.Println(len(tx), dst)
	byteArr := [][]byte{iv, tx}
	return bytes.Join(byteArr, []byte("")), nil
}

func AesGcmDecrypt(src []byte, key string) ([]byte, error) {
	cipherText := src
	//length := len(cipherText)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	iv := cipherText[0:12]
	//tag := cipherText[length-16:]
	//data := cipherText[12 : length-16]
	data := cipherText[12:]
	//tag := nil
	plaintext, err := aesgcm.Open(nil, iv, data, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
