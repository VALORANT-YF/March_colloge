package usersModel

// UserLoginResp 用户登录之后,需要返回的信息
type UserLoginResp struct {
	AToken string `json:"aToken"`
	IsBoss bool   `json:"is_boss"`
}
