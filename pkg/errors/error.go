package myerrors

import "errors"

var (
	ErrCodeParamInvalid = 4001 //参数错误
	ErrCodeUserExisted  = 4002 //用户已存在
	ErrCodeRegisterFail = 4003 //注册失败
)

var (
	ErrUserExisted = errors.New("user existed")
)
