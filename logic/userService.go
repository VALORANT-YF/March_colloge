package logic

import (
	"college/dao/mysql"
	"college/models/bookBlogArticle"
	"college/models/usersModel"
	"college/pkg/jwtToken"
	"strconv"
	"strings"
)

// FindUserIfExistService 查询用户信息,如果用户存在,则允许登录
func FindUserIfExistService(userAccount *usersModel.UserAccount) (error, bool, string, usersModel.UserLoginResp) {
	mobile := userAccount.Mobile     //电话号码
	name := userAccount.Name         //用户姓名
	password := userAccount.Password //密码

	var aToken string
	userInformation, err := mysql.SelectUserByTelAndName(mobile, name, password)
	var resp = usersModel.UserLoginResp{}
	//如果发生错误或者用户的唯一表示unionid为空,则登录失败
	if err != nil || len(userInformation.Unionid) == 0 {
		return err, false, aToken, resp
	} else {
		//如果用户登录成功
		aToken, _, _ = jwtToken.GenToken(userInformation.Unionid) //生成aToken
	}
	resp.IsBoss = userInformation.IsBoss
	resp.AToken = aToken
	return nil, true, aToken, resp
}

// FindBookBlogAddressService 根据unionid查询用户简书博客的主页链接 判断用户是否是第一次登录
func FindBookBlogAddressService(unionid string) (error, bool) {
	err, address := mysql.SelectBookBlogAddress(unionid)
	if err != nil {
		return err, false
	}
	if len(address.BlogAddress) == 0 {
		return nil, false
	}
	return nil, true
}

// UpdateBookBlogAddressService 新用户更新简书和博客的主页地址
func UpdateBookBlogAddressService(unionid string, address usersModel.BlogBookAddress) error {
	var tbUserAddress = new(usersModel.TbUser)
	tbUserAddress.BlogAddress = address.BlogAddress   //博客主页地址
	tbUserAddress.BooksAddress = address.BooksAddress //简书主页地址
	err := mysql.UpdateBookBlogAddress(unionid, tbUserAddress)
	if err != nil {
		return err
	}
	return nil
}

// GetBookOrBlogService 用户查看公共区域简书或者博客的文章
func GetBookOrBlogService(userLook usersModel.UserLookBookOrBlog) (err error, articleDetails usersModel.UserLookBookBlogDetails) {
	id := userLook.Id                     //简书或者博客文章的id
	isBookOrBlog := userLook.IsBookOrBlog //区分是简书还是博客 0简书 1博客
	if isBookOrBlog == 0 {
		var book bookBlogArticle.TbBookArticle
		err, book = mysql.SelectBookById(int64(id))
		articleDetails.Id = int64(book.ID)            //文章id
		articleDetails.BookTitle = book.BookTitle     //文章标题
		articleDetails.BookArticle = book.BookArticle //文章详细内容
		articleDetails.Name = book.Name               //文章主人的姓名
		articleDetails.DeptName = book.DeptName       //文章主人所在的部门
		articleDetails.Url = book.ArticleUrl          //原文链接
	} else {
		var blog bookBlogArticle.TbBlog
		err, blog = mysql.SelectBlogById(int64(id))
		articleDetails.Id = int64(blog.ID)            //文章id
		articleDetails.BookTitle = blog.BookTitle     //文章标题
		articleDetails.BookArticle = blog.BookArticle //文章详细内容
		articleDetails.Name = blog.Name               //文章主人的姓名
		articleDetails.DeptName = blog.DeptName       //文章主人所在的部门
		articleDetails.Url = blog.BlogUrl             //原文链接
	}
	return
}

// GetUserSelfInformationService 得到当前登录用户的个人信息
func GetUserSelfInformationService(unionid string) (err error, selfUserInformation usersModel.UserSelfInformation) {
	err, user := mysql.SelectSelfInformation(unionid)
	//根据部门编号拿到部门名称
	deptList := user.DeptIdList
	// 去掉字符串中的方括号和空格
	deptList = strings.Trim(deptList, "[]")
	// 分割字符串并存储到切片
	parts := strings.Split(deptList, " ")
	// 创建一个整数切片，用于存储解析后的整数
	var numbers []int64
	// 遍历字符串切片并将每个元素解析为整数
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err == nil {
			numbers = append(numbers, int64(num))
		}
	}
	//根据部门编号拿到部门名称集合
	err, deptName := mysql.SelectDeptName(numbers)
	if err != nil {
		return
	}
	deptNameStr := strings.Join(deptName, " ")
	selfUserInformation.Name = user.Name                 //姓名
	selfUserInformation.Avatar = user.Avatar             //头像地址
	selfUserInformation.BlogAddress = user.BlogAddress   //博客地址
	selfUserInformation.BooksAddress = user.BooksAddress //简书地址
	selfUserInformation.Mobile = user.Mobile             //电话
	selfUserInformation.DeptName = deptNameStr           //部门名称
	selfUserInformation.Password = user.Password         //密码
	return
}

// UpdateSelfUserInformationService 更新用户个人信息
func UpdateSelfUserInformationService(unionid string, information usersModel.UserUpdateSelfInformation) (error, uint8) {
	var isUpdatePassword uint8
	isUpdatePassword = 1
	//拿到用户个人信息
	err, u := mysql.SelectSelfInformation(unionid)
	if err != nil {
		return err, isUpdatePassword //1表示密码没有被修改
	}
	if u.Password != information.Password {
		isUpdatePassword = 0 //0表示密码被修改
	}
	return mysql.UpdateUserSelfInformationDao(unionid, information), isUpdatePassword
}
