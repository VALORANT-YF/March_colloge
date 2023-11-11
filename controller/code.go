package controller

type ResCode int64

const CtxtUserIDKey = "userID"

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CondeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeValidToken
	CodeUpdatePassword
	CodeNoAdmin
	CodeNoURL
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:          "success",
	CodeInvalidParam:     "请求参数错误",
	CodeUserExist:        "用户名已经存在",
	CodeUserNotExist:     "用户名不存在",
	CondeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:       "服务繁忙",
	CodeNeedLogin:        "用户未登录",
	CodeValidToken:       "登录错误",
	CodeUpdatePassword:   "密码被修改",
	CodeNoAdmin:          "权限不足",
	CodeNoURL:            "简书博客链接为空",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
