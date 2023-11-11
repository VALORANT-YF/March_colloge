package dingToken

import (
	"college/logic/dingOfficialService"
	"go.uber.org/zap"
)

// GetOfficialAccessToken	从Redis数据库中拿到accessToken
func GetOfficialAccessToken() string {
	accessToken, err := dingOfficialService.ObtainTokenService()
	if err != nil {
		zap.L().Error("dingOfficialService.ObtainTokenService() is failed", zap.Error(err))
		return ""
	}
	return accessToken
}
