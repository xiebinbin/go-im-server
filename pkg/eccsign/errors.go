package eccsign

import "errors"

var (
	ErrSdkAppIDNotMatch    = errors.New("skey not match")
	ErrIdentifierNotMatch  = errors.New("identifier not match")
	ErrExpired             = errors.New("expired")
	ErrSigNotMatch         = errors.New("sig not match")
)
