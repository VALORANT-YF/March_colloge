package md5Password

import (
	"crypto/md5"
	"encoding/hex"
)

//对密码进行简单加密

const sercret = "zyf"

// EncryptPassword 对密码进行加密
func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(sercret))
	h.Sum([]byte(oPassword))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
