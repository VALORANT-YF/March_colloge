package controller

import (
	"college/logic"
	"college/models/usersModel"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct{}

// UserLogin 用户登录
func (u UserController) UserLogin(context *gin.Context) {
	//获取请求参数
	userAccount := new(usersModel.UserAccount)
	err := context.ShouldBindJSON(&userAccount)
	if err != nil {
		zap.L().Error("context.ShouldBindJSON(&userAccount) is failed", zap.Error(err))
		ResponseError(context, CodeInvalidParam)
		return
	}
	//判断账号是否存在
	err, isExist, _, resp := logic.FindUserIfExistService(userAccount)
	//如果用户账号不存在
	if !isExist {
		ResponseError(context, CodeUserNotExist)
		return
	}
	//返回响应参数
	ResponseSuccessWithData(context, resp)
}

// UserInitialBookBlog 用户第一次登录时需要填写自己的简书博客链接,即可以根据简书和博客地址是否为空判断用户是否是第一次登录
func (u UserController) UserInitialBookBlog(context *gin.Context) {
	//首先获得用户的unionid (存储在Token中)
	unionid := GetUnionidToken(context)
	if unionid == "no" {
		ResponseError(context, CodeNeedLogin)
		return
	}
	fmt.Println(unionid)
	err, bookisExists := logic.FindBookBlogAddressService(unionid)
	if err != nil {
		zap.L().Error("logic.FindBookBlogAddressService(unionid) is failed", zap.Error(err))
		return
	}
	//如果简书博客主页地址不存在 , 说明用户是第一次登录
	if !bookisExists {
		fmt.Println("简书链接不存在")
		fmt.Println(CodeNoURL)
		ResponseError(context, CodeNoURL)
		return
	}
	ResponseSuccess(context)
}

// UserInputBookBlogAddress 新用户录入自己的简书博客链接
func (u UserController) UserInputBookBlogAddress(context *gin.Context) {
	//获得用户的唯一id unionid
	unionid := GetUnionidToken(context)
	if unionid == "no" {
		ResponseError(context, CodeNeedLogin)
		return
	}
	//获得请求参数 即简书博客主页地址
	var address usersModel.BlogBookAddress
	err := context.ShouldBindJSON(&address)
	if err != nil {
		zap.L().Error("context.ShouldBindJSON(&address) is failed", zap.Error(err))
		ResponseError(context, CodeInvalidParam)
		return
	}
	//调用service层处理数据
	err = logic.UpdateBookBlogAddressService(unionid, address)
	if err != nil {
		zap.L().Error("logic.UpdateBookBlogAddressService(unionid, address) is failed", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	ResponseSuccess(context)
}

// GetOtherBookBlog 用户查看简书或者博客的内容
func (u UserController) GetOtherBookBlog(context *gin.Context) {
	//验证用户登录状态
	unionid := GetUnionidToken(context)
	if unionid == "no" {
		ResponseError(context, CodeNeedLogin)
		return
	}
	//接收请求参数
	var userLook usersModel.UserLookBookOrBlog
	if err := context.ShouldBindQuery(&userLook); err != nil {
		zap.L().Error("context.ShouldBindJSON(&userLook) is failed", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	err, userArticleOrBlogContent := logic.GetBookOrBlogService(userLook)
	if len(userArticleOrBlogContent.BookTitle) == 0 || err != nil {
		ResponseError(context, CodeServerBusy)
		return
	}
	ResponseSuccessWithData(context, userArticleOrBlogContent)
}

// GetSelfInformation 用户登录之后,拿到自己的个人信息
func (u UserController) GetSelfInformation(context *gin.Context) {
	//拿到当前登录用户的唯一id
	uniond := GetUnionidToken(context)
	if uniond == "no" {
		ResponseError(context, CodeNeedLogin)
		return
	}
	err, userSelfInformation := logic.GetUserSelfInformationService(uniond)
	if err != nil {
		zap.L().Error("logic.GetUserSelfInformationService(uniond) is failed", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	ResponseSuccessWithData(context, userSelfInformation)
}

// UpdateSelfInformation 用户修改个人信息
func (u UserController) UpdateSelfInformation(context *gin.Context) {
	//拿到当前登录用户的唯一id
	unionid := GetUnionidToken(context)
	if unionid == "no" {
		ResponseError(context, CodeServerBusy)
		return
	}
	//接收 请求参数
	var updateSelfInformation usersModel.UserUpdateSelfInformation
	err := context.ShouldBindJSON(&updateSelfInformation)
	if err != nil {
		ResponseError(context, CodeInvalidParam)
		return
	}
	var isUpdatePassword uint8
	if err, isUpdatePassword = logic.UpdateSelfUserInformationService(unionid, updateSelfInformation); err != nil {
		ResponseError(context, CodeServerBusy)
		return
	}
	//如果密码被修改
	if isUpdatePassword == 0 {
		ResponseSuccessWithData(context, isUpdatePassword)
	} else {
		ResponseSuccessWithData(context, CodeSuccess)
	}
}
