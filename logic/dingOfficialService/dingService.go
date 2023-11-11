package dingOfficialService

import (
	"college/dao/mysql/dingMysql"
	"college/dao/redis/dingRedisDao"
	"college/models/deptsModel"
	"college/models/dingOfficialModel"
	"college/models/usersModel"
	"fmt"
)

// GetTokenService 将得到的Token存到Redis中
func GetTokenService(accessToken string) error {
	return dingRedisDao.GetAccessTokenDao(accessToken) //调用redis中的操作将accessToken存储在redis数据库中
}

// ObtainTokenService 获得Redis数据库中存储的Token
func ObtainTokenService() (string, error) {
	return dingRedisDao.ObtainTokenDao()
}

// CreateDeptInformationService 将得到的部门信息插入数据库中
func CreateDeptInformationService(deptSearchResult []deptsModel.TbDept) error {
	//保留已经被修改的部门的信息
	for i, dept := range deptSearchResult {
		err, d := dingMysql.SelectDeptExistsInformation(dept.DeptId)
		if err != nil {
			return err
		}
		if d.Name != "" {
			deptSearchResult[i].IsWriteBooks = d.IsWriteBooks //部门是否需要写简书
		}
	}

	//每一次插入先删除
	if err := dingMysql.DeleteAllDept(); err != nil {
		return err
	}

	//只去插入不存在的部门
	return dingMysql.InsertDepts(deptSearchResult)
}

// SelectDeptIdsService 获取部门id
func SelectDeptIdsService() (err error, deptIds []int64) {
	deptIds, err = dingMysql.SelectRightDept() //拿到所有的部门id
	return
}

// CreateUserInformationService 向tb_user表中插入数据
func CreateUserInformationService(usersInformationResult []dingOfficialModel.UsersInformationResult) error {

	var tbUsers []usersModel.TbUser
	for i := 0; i < len(usersInformationResult); i++ {
		var tbUser usersModel.TbUser //每一次循环创建一个新的结构
		//给tbUser结构体赋值
		for j := 0; j < len(usersInformationResult[i].Result.List); j++ {
			temp := usersInformationResult[i].Result.List[j] //temp中间变量
			tbUser.DeptIdList = fmt.Sprintf("%v", temp.DeptIdList)
			tbUser.Name = temp.Name
			tbUser.IsBoss = temp.Boss
			tbUser.Userid = temp.Userid
			tbUser.Unionid = temp.Unionid
			tbUser.Avatar = temp.Avatar
			tbUser.Mobile = temp.Mobile
			tbUser.Password = temp.Mobile[:6] //设置一个初始密码
			tbUsers = append(tbUsers, tbUser) //将全部结果拼接
		}
	}
	rightTbUsers := make([]usersModel.TbUser, 0, len(tbUsers))
	//删掉重复出现的人
	encountered := map[usersModel.TbUser]bool{}
	//遍历tbUsers切片
	for v := range tbUsers {
		if encountered[tbUsers[v]] == true {
			// 如果元素已经在map中存在，跳过
			continue
		} else {
			// 如果元素不在map中，将其添加到结果切片中，并在map中标记为已经遇到
			encountered[tbUsers[v]] = true
			rightTbUsers = append(rightTbUsers, tbUsers[v])
		}
	}

	//保留用户修改的数据
	for i, user := range rightTbUsers {
		err, u := dingMysql.SelectUserExistInformation(user.Unionid)
		if err != nil {
			return err
		}
		if u.Name != "" {
			// 使用指针来更新切片元素
			rightTbUsers[i].IsBoss = u.IsBoss                   //用户是否是管理员
			rightTbUsers[i].NotWrittenCount = u.NotWrittenCount //未写次数
			rightTbUsers[i].BlogAddress = u.BlogAddress         //博客地址
			rightTbUsers[i].BooksAddress = u.BooksAddress       //简书地址
			rightTbUsers[i].Password = u.Password               //用户密码
			rightTbUsers[i].ExcellentCount = u.ExcellentCount   //优秀次数
		}
	}

	//每一次插入先删除
	if err := dingMysql.DeleteAllUser(); err != nil {
		return err
	}

	//调用dao层插入数据
	err := dingMysql.InsertUsersInformation(rightTbUsers)
	if err != nil {
		return err
	}
	return nil
}
