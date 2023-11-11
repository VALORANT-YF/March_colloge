package dingRedisDao

import "college/dao/redis"

// GetAccessTokenDao 调用官方接口获取AccessToken , 并且将AccessToken存储在Redis中
func GetAccessTokenDao(accessToken string) error {
	return redis.RDB.Set("accessToken", accessToken, 0).Err() //将企业内部应用的AccessToken存储在数据库中,0表示该Token永不过期
}

// ObtainTokenDao 获得从Redis数据库中存储的Token
func ObtainTokenDao() (string, error) {
	accessToken, err := redis.RDB.Get("accessToken").Result()
	return accessToken, err
}
