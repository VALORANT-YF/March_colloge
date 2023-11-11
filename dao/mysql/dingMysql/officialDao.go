package dingMysql

import (
	"college/dao/mysql"
	"college/models/deptsModel"
	"college/models/usersModel"

	"go.uber.org/zap"
)

// InsertDepts 将查询出的部门批量插入数据库总
func InsertDepts(deptSearchResult []deptsModel.TbDept) (err error) {
	//开启事务
	tx := mysql.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() //发生错误回滚事务
		}
	}()

	if tx.Error != nil {
		zap.L().Error("tx.Error", zap.Error(err))
		return err
	}

	//在事务中指定插入操作
	for _, value := range deptSearchResult {
		err := tx.Create(&value).Error
		if err != nil {
			//回滚事务
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error //提交事务,如果有错误,返回错误
}

// SelectRightDept 查询出部门id
func SelectRightDept() ([]int64, error) {
	var deptIds []int64
	result := mysql.DB.Model(&deptsModel.TbDept{}).Pluck("dept_id", &deptIds)
	return deptIds, result.Error
}

// InsertUsersInformation 向tb_user表中插入用户信息
func InsertUsersInformation(tbUsers []usersModel.TbUser) error {
	//开启事务
	tx := mysql.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//在事务中指定插入操作
	for _, value := range tbUsers {
		err := tx.Create(&value).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// DeleteAllUser 清空user表
func DeleteAllUser() error {
	return mysql.DB.Exec("TRUNCATE TABLE tb_user").Error
}

// DeleteAllDept 清空dept表
func DeleteAllDept() error {
	return mysql.DB.Unscoped().Exec("TRUNCATE TABLE tb_dept").Error
}

// SelectUserExistInformation 查询用户已经存在的信息
func SelectUserExistInformation(unionid string) (err error, userResult usersModel.TbUser) {
	err = mysql.DB.Where("unionid = ?", unionid).Find(&userResult).Error
	return
}

// SelectDeptExistsInformation 查询部门已经存在的信息
func SelectDeptExistsInformation(deptId int64) (err error, deptResult deptsModel.TbDept) {
	err = mysql.DB.Where("dept_id = ?", deptId).Find(&deptResult).Error
	return
}
