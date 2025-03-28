package myerrors

import "errors"

var (
	ErrCodeParamInvalid   = 4001 //参数错误
	ErrCodeUserExisted    = 4002 //用户已存在
	ErrCodeRegisterFail   = 4003 //注册失败
	ErrCodeMissingToken   = 4004 //未携带token
	ErrCodeTokenExpired   = 4005 //token已失效
	ErrCodeInvalidToken   = 4006 //token错误
	ErrCodeUserNotExisted = 4007 //用户不存在
	ErrCodePasswordWrong  = 4008 //密码错误
	ErrCodeLoginFail      = 4009 //登陆失败
)

var (
	ErrUserExisted    = errors.New("user existed")
	ErrPasswordWrong  = errors.New("password wrong")
	ErrTokenExpired   = errors.New("token expired")
	ErrUserNotExisted = errors.New("user not existed")
)
