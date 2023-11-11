package mysql

import (
	"college/models/robotModels"
)

// SelectRobotToken 管理员查询机器人的Token
func SelectRobotToken() (err error, result []robotModels.TbRobot) {
	err = DB.Find(&result).Error
	return
}

// InsertRobotToken 管理员插入新的机器人token
func InsertRobotToken(newToken robotModels.TbRobot) (err error) {
	return DB.Create(&newToken).Error
}

// UpdateRobotToken 管理员修改机器人的token
func UpdateRobotToken(newToken robotModels.TbRobot) (err error) {
	return DB.Updates(&newToken).Error
}

// DeleteRobotToken 管理员删除Token
func DeleteRobotToken(token robotModels.TbRobot) (err error) {
	return DB.Delete(&token).Error
}
