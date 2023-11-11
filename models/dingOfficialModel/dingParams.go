package dingOfficialModel

// GetTokenParams 获取企业内部应用的accessToken需要传递的参数
type GetTokenParams struct {
	Appkey    string `json:"appkey,omitempty" binding:"required"`    //应用的唯一标识key
	Appsecret string `json:"appsecret,omitempty" binding:"required"` //应用的密钥
}
