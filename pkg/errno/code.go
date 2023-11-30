package errno

const (
	OK           = 200
	DefErr       = 400
	TokenErr     = 401
	Exception    = 402
	WrongReq     = 422
	SysErr       = 500
	HeaderErr    = 510
	TimesLimited = 512

	SaveDataFailed  = 100000
	ParamsErr       = 100003
	ParamsFormatErr = 100004
	MissingParams   = 100005
	DataNotExist    = 100006
	FORBIDDEN       = 100007
	DataMustUnique  = 100008
	UserNotExist    = 100010
	UserDelete      = 100011
	UserUnavailable = 100012
	DataDeleted     = 100013
	EnvNotMatch     = 100014

	CaptchaUnavailable = 400017

	Expired         = 200001
	QrCodeFormatErr = 200002
	QrCodeInvalid   = 200003
	UserNotInChat   = 200004
)
