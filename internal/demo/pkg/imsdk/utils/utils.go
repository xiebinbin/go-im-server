package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func GenerateNonce() string {
	u4 := uuid.New()
	return u4.String()
}

func Hash256(str []byte) []byte {
	h := sha256.New()
	h.Write(str)
	return h.Sum(nil)
}

func GetTimeSecs() int64 {
	return time.Now().Unix()
}

func GetRandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result[i] = bytes[r.Intn(len(bytes))]
	}
	return string(result)
}

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Md516(str string) string {
	res := Md5Str(str)
	return res[8:24]
}

func GetNanos() int64 {
	return time.Now().UnixNano()
}

func GetMillis() int64 {
	return GetNanos() / 1e6
}
