package mysql

import (
	"college/models/deptsModel"
	"college/models/usersModel"
	"gorm.io/gorm"
)

const (
	tableName = "tb_user"
)

// SelectUserByTelAndName 根据电话以及用户姓名查询用户是否存在
func SelectUserByTelAndName(mobile, name, password string) (*usersModel.TbUser, error) {
	userInformation := new(usersModel.TbUser)
	err := DB.Where("name = ? and mobile = ? and password = ?", name, mobile, password).Find(&userInformation).Error
	return userInformation, err
}

// SelectBookBlogAddress 根据unionid查询用户简书和博客的主页地址
func SelectBookBlogAddress(unionid string) (error, usersModel.TbUserBookBlog) {
	var bookBlogAddress usersModel.TbUserBookBlog
	err := DB.Where("unionid = ?", unionid).Table("tb_user").Find(&bookBlogAddress).Error
	return err, bookBlogAddress
}

// UpdateBookBlogAddress 用户插入简书博客主页地址 (使用更新操作即可,并不是新插入一条数据)
func UpdateBookBlogAddress(unionid string, address *usersModel.TbUser) error {
	err := DB.Where("unionid = ?", unionid).Updates(&address).Error
	return err
}

// SelectDeptListName 根据unionid 去查询用户的部门列表和用户姓名
func SelectDeptListName(unionid string) (error, usersModel.TbUserDeptAndName) {
	var userDeptAndName usersModel.TbUserDeptAndName
	err := DB.Where("unionid = ?", unionid).Table("tb_user").Find(&userDeptAndName).Error
	return err, userDeptAndName
}

// SelectDeptName 根据部门编号查询部门名称
func SelectDeptName(deptIds []int64) (error, []string) {
	var deptName []string
	var tbDeptName deptsModel.TbDeptName
	for _, deptId := range deptIds {
		err := DB.Where("dept_id = ?", deptId).Table("tb_dept").Find(&tbDeptName).Error
		if err != nil {
			return err, nil
		}
		deptName = append(deptName, tbDeptName.Name)
	}
	return nil, deptName
}

// SelectBookAddressPerson 查询每一个人的简书博客地址,电话,姓名,以及部门编号
func SelectBookAddressPerson() (error, []usersModel.AddressPerson) {
	var result []usersModel.AddressPerson
	err := DB.Select("name , books_address , blog_address , dept_id_list , mobile").Find(&result).Error
	return err, result
}

// SelectAllPersonInformation 查询部门中所有人的信息
func SelectAllPersonInformation(name string) (error, []usersModel.TbUser) {
	var allPersonInformation []usersModel.TbUser
	err := DB.Select("name , dept_id_list , mobile , books_address , blog_address , userid , is_boss").Where("name like ?", "%"+name+"%").Find(&allPersonInformation).Error
	return err, allPersonInformation
}

// SelectSelfInformation 查看自己的信息
func SelectSelfInformation(unionid string) (err error, userSelfInformation usersModel.TbUser) {
	err = DB.Select("is_boss , excellent_count ,not_written_count , mobile , avatar , name , books_address , blog_address , dept_id_list , password").Where("unionid = ?", unionid).Find(&userSelfInformation).Error
	return
}

// UpdateUserSelfInformationDao 修改用户个人信息
func UpdateUserSelfInformationDao(unionid string, information usersModel.UserUpdateSelfInformation) (err error) {
	err = DB.Where("unionid = ?", unionid).Table("tb_user").Updates(&information).Error
	return
}

// UpdateUserExcellentCount 更新用户优秀简书和博客的次数
/*
id : 用户id
newExcellentCount : 新的优秀简书的次数
*/
func UpdateUserExcellentCount(id uint, newExcellentCount uint32, tx *gorm.DB) error {
	return tx.Select("excellent_count").Table(tableName).Where("id = ?", id).Update("excellent_count", newExcellentCount).Error
}

// SelectInformationByTel 根据电话查找用户信息
func SelectInformationByTel(mobile string) (error, usersModel.TbUser) {
	var result usersModel.TbUser
	err := DB.Select("id , is_boss , excellent_count , not_written_count , avatar , name").Where("mobile = ?", mobile).Find(&result).Error
	return err, result
}

// SelectMobileByDeptId 根据部门id 去查询该部门下面所有人的信息
func SelectMobileByDeptId(deptId string) (err error, users []usersModel.TbUser) {
	err = DB.Where("CONCAT('[',dept_id_list,']') LIKE ?", "%"+deptId+"%").Find(&users).Error
	return
}

// UpdateUserNoWriteCount 更新用户简书未写次数
func UpdateUserNoWriteCount(userID string, newCount uint32) error {
	return DB.Where("userid = ?", userID).Table(tableName).Update("not_written_count", newCount).Error
}
