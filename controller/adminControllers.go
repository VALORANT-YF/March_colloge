package controller

import (
	"college/logic"
	"college/logic/dingOfficialService"
	"college/models/bookBlogArticle"
	"college/models/dingOfficialModel"
	"college/models/robotModels"
	"college/models/usersModel"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type AdminController struct{}

// GetTokenController 获取企业内部应用的accessToken并存入Redis数据库中 (超级管理员登录)
func (d AdminController) GetTokenController(context *gin.Context) {
	//获取请求参数 appkey 和appsecret
	var getTokenParams = new(dingOfficialModel.GetTokenParams)
	err := context.ShouldBindJSON(&getTokenParams)
	if err != nil {
		zap.L().Error("context.ShouldBindJSON(&getTokenParams) is failed", zap.Error(err))
		return
	}
	//拿到请求参数,调用官方 获取accessToken
	appkey := getTokenParams.Appkey
	appSecret := getTokenParams.Appsecret
	//url := "https://oapi.dingtalk.com/gettoken?appkey=" + appKey + "&appsecret=" + appSecret
	//闭包 函数+引用环境
	accessToken := func() string {
		url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey="+"%v"+"&appsecret="+"%v", appkey, appSecret)
		//发送Get请求
		response, err := http.Get(url)
		if err != nil {
			zap.L().Error("http.Get(\"https://oapi.dingtalk.com/gettoken?appkey=\" + appkey + \"&appsecret=\" + appSecret) is failed", zap.Error(err))
			return ""
		}
		defer response.Body.Close()
		//读取响应主体
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			zap.L().Error("io.ReadAll(response.Body) is error", zap.Error(err))
			return ""
		}
		// 使用 gjson 库从响应中提取访问令牌
		return gjson.Get(string(responseBody), "access_token").String()
	}()
	if len(accessToken) == 0 { //如果调官方接口拿到的Token为空
		ResponseError(context, CondeInvalidPassword)
		return
	}
	err = dingOfficialService.GetTokenService(accessToken) //将得到的AccessToken存储在Redis数据库中
	if err != nil {
		zap.L().Error("dingOfficialService.GetTokenService(accessToken) is failed", zap.Error(err))
		return
	}
	ResponseSuccess(context)
}

// DeptIsWrite 设置部门是否需要写简书
func (d AdminController) DeptIsWrite(context *gin.Context) {
	//接收请求参数
	deptId := context.Param("dept_id")
	isWriteBook := context.Param("is_write_books")
	//参数校验
	if len(deptId) == 0 || (isWriteBook != "1" && isWriteBook != "0") {
		ResponseError(context, CodeInvalidParam)
		return
	}
	err := logic.UpdateDeptIsWriteService(deptId, isWriteBook)
	if err != nil {
		zap.L().Error("logic.UpdateDeptIsWriteService(deptId, isWriteBook) is error", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	ResponseSuccess(context)
}

// SelectDeptPersonInformation 查询出正确的部门架构和人员信息
func (d AdminController) SelectDeptPersonInformation(context *gin.Context) {
	//接收请求参数
	name := context.Query("name")
	err, deptUserInformation := logic.SelectDeptPersonInformationService(name)
	if err != nil {
		zap.L().Error("logic.SelectDeptPersonInformationService() is failed", zap.Error(err))
		ResponseError(context, CodeServerBusy)
	}
	ResponseSuccessWithData(context, deptUserInformation)
}

// UpdateAdmin 设置或取消管理员
func (d AdminController) UpdateAdmin(context *gin.Context) {
	//获取请求参数
	var adminUri usersModel.UpdateAdminUri
	if err := context.ShouldBindUri(&adminUri); err != nil {
		ResponseError(context, CodeInvalidParam)
		return
	}
	if adminUri.IsBoss != 0 && adminUri.IsBoss != 1 {
		ResponseError(context, CodeInvalidParam)
		return
	}
	//fmt.Println(adminUri)
	err := logic.UpdateAdminService(adminUri)
	if err != nil {
		zap.L().Error("logic.UpdateAdminService(adminUri) is failed", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	ResponseSuccess(context)
}

// ExcellentBookBlog 管理员评价优秀简书或者博客
func (d AdminController) ExcellentBookBlog(context *gin.Context) {
	//验证管理员登录状态 根据unionid 查看用户是否是管理员
	unionid := GetUnionidToken(context)
	if !GetIsBoss(unionid) {
		ResponseError(context, CodeNoAdmin)
		return
	}
	//获取请求参数
	var excellent bookBlogArticle.Excellent
	if err := context.ShouldBindJSON(&excellent); err != nil {
		ResponseError(context, CodeInvalidParam)
		return
	}
	if err := logic.ElectExcellentUserService(excellent); err != nil {
		zap.L().Error("logic.ElectExcellentUserService(excellent) is failed", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	ResponseSuccess(context)
}

// LookRobot 管理员查看机器人token
func (d AdminController) LookRobot(context *gin.Context) {
	//验证管理员登录状态 根据unionid 查看用户是否是管理员
	unionid := GetUnionidToken(context)
	if !GetIsBoss(unionid) {
		ResponseError(context, CodeNoAdmin)
		return
	}
	err, result := logic.GetRobotTokenList()
	if err != nil {
		ResponseError(context, CodeServerBusy)
		return
	}
	ResponseSuccessWithData(context, result)
}

// AddRobot 管理员新增机器人
func (d AdminController) AddRobot(context *gin.Context) {
	//验证管理员登录状态 根据unionid 查看用户是否是管理员
	unionid := GetUnionidToken(context)
	if !GetIsBoss(unionid) {
		ResponseError(context, CodeNoAdmin)
		return
	}
	//获取请求参数
	newRobot := robotModels.TbRobot{}
	err := context.ShouldBindJSON(&newRobot)
	if err != nil {
		ResponseError(context, CodeInvalidParam)
		return
	}
	err = logic.AddRobotToken(unionid, newRobot)
	if err != nil {
		zap.L().Error("logic.AddRobotToken(unionid , newRobot) is failed", zap.Error(err))
		return
	}
	ResponseSuccess(context)
}

// UpdateRobot 管理员修改机器人token
func (d AdminController) UpdateRobot(context *gin.Context) {
	//验证管理员登录状态 根据unionid 查看用户是否是管理员
	unionid := GetUnionidToken(context)
	if !GetIsBoss(unionid) {
		ResponseError(context, CodeNoAdmin)
		return
	}
	//获取请求参数
	newRobot := robotModels.TbRobot{}
	err := context.ShouldBindJSON(&newRobot)
	if err != nil {
		ResponseError(context, CodeInvalidParam)
		return
	}
	err = logic.ChangeRobotToken(unionid, newRobot)
	if err != nil {
		zap.L().Error("logic.AddRobotToken(unionid , newRobot) is failed", zap.Error(err))
		return
	}
	ResponseSuccess(context)
}

// DeleteRobot 管理员删除机器人
func (d AdminController) DeleteRobot(context *gin.Context) {
	//验证管理员登录状态 根据unionid 查看用户是否是管理员
	unionid := GetUnionidToken(context)
	if !GetIsBoss(unionid) {
		ResponseError(context, CodeNoAdmin)
		return
	}
	//获取请求参数
	newRobot := robotModels.TbRobot{}
	err := context.ShouldBindJSON(&newRobot)
	if err != nil {
		ResponseError(context, CodeInvalidParam)
		return
	}
	err = logic.DropRobotToken(unionid, newRobot)
	if err != nil {
		zap.L().Error("logic.AddRobotToken(unionid , newRobot) is failed", zap.Error(err))
		return
	}
	ResponseSuccess(context)
}

/*
下面是定时任务 测试使用
*/

// NoWriteUsers 未写简书 博客的人 的名单
func (d AdminController) NoWriteUsers(context *gin.Context) {
	//验证管理员登录状态 鉴权
	//unionid := GetUnionidToken(context)
	//if !GetIsBoss(unionid) {
	//	ResponseError(context, CodeNoAdmin)
	//	return
	//}
	err, noWriteList := logic.SelectNoWriteUserService()
	if err != nil {
		zap.L().Error("logic.SelectNoWriteUserService() is error", zap.Error(err))
		ResponseError(context, CodeServerBusy)
	}
	ResponseSuccessWithData(context, noWriteList)
}

// ExcellentPersonList 查询出本周的优秀简书博客名单
func (d AdminController) ExcellentPersonList(context *gin.Context) {
	//验证管理员登录状态 鉴权
	//unionid := GetUnionidToken(context)
	//if !GetIsBoss(unionid) {
	//	ResponseError(context, CodeNoAdmin)
	//	return
	//}
	ResponseSuccessWithData(context, logic.SelectExcellentPersonService())
}

// ExcellentCountFive 查询出优秀简书次数前5的人
func (d AdminController) ExcellentCountFive(context *gin.Context) {
	err, result := logic.SelectExcellentCountService()
	if err != nil {
		zap.L().Error("logic.SelectExcellentCountService()", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	ResponseSuccessWithData(context, result)
}

// NoWriteCountThree 查询未写次数前3的人
func (d AdminController) NoWriteCountThree(context *gin.Context) {
	err, result := logic.SelectNoWriteCount()
	if err != nil {
		zap.L().Error("logic.SelectNoWriteUserService()", zap.Error(err))
		ResponseError(context, CodeServerBusy)
		return
	}
	ResponseSuccessWithData(context, result)
}
