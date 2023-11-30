package errors

import "imsdk/pkg/errno"

var (
	ErrSdkDefErr           = errno.Add("default error", 400)
	ErrSdkDataNotExists    = errno.Add("data not exists", 100000)
	ErrSdkTimeOut          = errno.Add("auth time expire", 100001)
	ErrSdkRequestRepeat    = errno.Add("request repeat", 100002)
	ErrSdkUserRegisterFail = errno.Add("user register fail", 100003)
	ErrSdkAKNotMatch       = errno.Add("auth verify fail", 100100)
	ErrSdkSignVerifyFail   = errno.Add("sign verify fail", 100101)
)
