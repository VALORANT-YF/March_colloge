package logic

import (
	"college/dao/mysql"
	"college/models/bookBlogArticle"
	"college/models/robotModels"
	"college/models/usersModel"
	"errors"
	"strconv"
	"time"
)

// QueryIsBossService 对登录用户鉴权
func QueryIsBossService(unionid string) bool {
	_, u := mysql.SelectSelfInformation(unionid)
	return u.IsBoss
}

// UpdateDeptIsWriteService 修改部门是否需要填写简书
func UpdateDeptIsWriteService(deptIdStr string, isWrite string) error {
	var err error
	deptId, err := strconv.ParseInt(deptIdStr, 10, 64)
	if err != nil {
		return err
	}
	if isWrite == "1" {
		err = mysql.UpdateDeptIsWrite(1, deptId)
	} else {
		err = mysql.UpdateDeptIsWrite(0, deptId)
	}
	return err
}

// UpdateAdminService 设置或者取消管理员
func UpdateAdminService(adminUri usersModel.UpdateAdminUri) error {
	isBoss := adminUri.IsBoss
	userid := adminUri.Userid
	return mysql.UpdateAdmin(isBoss, userid)
}

// SelectNoWriteUserService 查询出简书博客未写的名单
func SelectNoWriteUserService() (err error, queryResult []usersModel.NoWriteView) {
	var noWriteList []usersModel.NoWriteUser //未写简书的人员名单
	//1.拿到需要写简书的部门
	err, needWriteDept := mysql.SelectNeedWriteDept()
	if err != nil {
		return
	}
	//2.根据需要写简书的部门的部门id , 去查找该部门下面的所有人的电话号码
	for _, data := range needWriteDept {
		deptId := strconv.FormatInt(data.DeptId, 10) //拿到部门id 转化为字符串
		//根据部门id , 查找电话号码
		err, users := mysql.SelectMobileByDeptId(deptId)
		if err != nil {
			return err, nil
		}
		for _, user := range users {
			//去根据电话号码查找简书 或者 博客 两个都没找到即为没写
			_, titleBook := mysql.SelectArticleByMobile(user.Mobile)
			_, titleBlog := mysql.SelectBlogByMobile(user.Mobile)
			if len(titleBook) == 0 && len(titleBlog) == 0 {
				//判断现在时间是否是周日晚上11:30 , 如果是未写次数加1
				now := time.Now()
				//计算目标时间 即 周日晚上11:30
				targetTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 30, 0, 0, now.Location())
				// 允许的时间误差为20分钟
				allowedTimeDelta := 20 * time.Minute
				if now.After(targetTime.Add(-allowedTimeDelta)) && now.Before(targetTime.Add(allowedTimeDelta)) && now.Weekday() == time.Sunday {
					// 当前时间在目标时间范围内，并且是周日晚上11:30附近
					err = mysql.UpdateUserNoWriteCount(user.Userid, user.NotWrittenCount+1)
					if err != nil {
						return err, nil
					}
				}
				//这个人没有写简书
				//fmt.Println("简书未写的人:" + user.Name)
				//根据未写简书的人的部门id 查询此人所在的部门
				//fmt.Println("此人的部门:" + data.Name)
				noWriteUser := usersModel.NoWriteUser{
					NotWrittenCount: user.NotWrittenCount, //简书未写的次数
					Mobile:          user.Mobile,          //未写简书的人的电话
					UserId:          user.Userid,          //未写简书的人的id
					Avatar:          user.Avatar,          //未写简书的人的头像
					Name:            user.Name,            //未写简书的人的名称
					BooksAddress:    user.BooksAddress,    //此人的简书主页地址
					BlogAddress:     user.BlogAddress,     //此人的博客主页地址
					DeptName:        data.Name,            //此人所在的部门
				}
				noWriteList = append(noWriteList, noWriteUser)
			}
		}
	}
	resultMap := make(map[string][]usersModel.NoWriteUser, len(noWriteList))
	for _, noWrite := range noWriteList {
		resultMap[noWrite.DeptName] = append(resultMap[noWrite.DeptName], noWrite)
	}
	var resultView []usersModel.NoWriteView
	for key, value := range resultMap {
		var end usersModel.NoWriteView
		end.OneDeptName = key
		end.NoWriteUserList = value
		resultView = append(resultView, end)
	}
	return nil, resultView
}

// SelectExcellentPersonService 查询本周的优秀简书博客
func SelectExcellentPersonService() []bookBlogArticle.ViewResult {
	//创建map集合
	resultMap := make(map[string][]interface{})
	_, excellentBook := mysql.SelectTopArticle()
	for _, value := range excellentBook {
		if _, ok := resultMap[value.DeptName]; !ok {
			resultMap[value.DeptName] = make([]interface{}, 0)
		}
		resultMap[value.DeptName] = append(resultMap[value.DeptName], value)
	}

	_, excellentBlog := mysql.SelectTopBlog()
	for _, value := range excellentBlog {
		if _, ok := resultMap[value.DeptName]; !ok {
			resultMap[value.DeptName] = make([]interface{}, 0)
		}
		resultMap[value.DeptName] = append(resultMap[value.DeptName], value)
	}
	var result []bookBlogArticle.ViewResult
	for key, value := range resultMap {
		temp := bookBlogArticle.ViewResult{
			DeptName: key,
			TypeName: value,
		}
		result = append(result, temp)
	}
	return result
}

// SelectExcellentCountService 查询优秀简书博客次数前五的人
func SelectExcellentCountService() (err error, result []usersModel.TbUser) {
	return mysql.SelectExcellentCount()
}

// SelectNoWriteCount 查询简书博客未写次数前3的人
func SelectNoWriteCount() (err error, result []usersModel.TbUser) {
	return mysql.SelectNoWriteCount()
}

// GetRobotTokenList 管理员查询机器人的token
func GetRobotTokenList() (err error, result []robotModels.TbRobot) {
	return mysql.SelectRobotToken()
}

// AddRobotToken 管理员新增机器人Token
func AddRobotToken(unionid string, newToken robotModels.TbRobot) (err error) {
	//根据unionid 查询登录的信息
	err, loginUserInformation := mysql.SelectSelfInformation(unionid)
	if err != nil {
		return
	}
	if loginUserInformation.Mobile != "17884712216" {
		return errors.New("权限不足")
	}
	//插入token
	return mysql.InsertRobotToken(newToken)
}

// ChangeRobotToken 管理员修改机器人Token
func ChangeRobotToken(unionid string, newToken robotModels.TbRobot) (err error) {
	//根据unionid 查询登录的信息
	err, loginUserInformation := mysql.SelectSelfInformation(unionid)
	if err != nil {
		return
	}
	if loginUserInformation.Mobile != "17884712216" {
		return errors.New("权限不足")
	}
	//修改token
	return mysql.UpdateRobotToken(newToken)
}

// DropRobotToken 管理员删除机器人Token
func DropRobotToken(unionid string, newToken robotModels.TbRobot) (err error) {
	//根据unionid 查询登录的信息
	err, loginUserInformation := mysql.SelectSelfInformation(unionid)
	if err != nil {
		return
	}
	if loginUserInformation.Mobile != "17884712216" {
		return errors.New("权限不足")
	}
	//删除token
	return mysql.DeleteRobotToken(newToken)
}
