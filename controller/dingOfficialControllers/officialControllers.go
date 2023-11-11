package dingOfficialControllers

import (
	"college/controller"
	"college/logic/dingOfficialService"
	"college/models/deptsModel"
	"college/models/dingOfficialModel"
	"college/pkg/dingToken"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OfficialController struct {
	Ding
}

// GetDeptList 获取全部的部门列表
func (o OfficialController) GetDeptList(context *gin.Context) {
	var deptSearchResult []deptsModel.TbDept //封装查询出的全部部门
	//从Redis数据库中拿到accessToken
	accessToken := dingToken.GetOfficialAccessToken()
	if len(accessToken) == 0 {
		context.JSON(http.StatusOK, "权限不足")
		return
	}
	//searchDeptLow查询最低一级的部门 start 闭包
	var searchDeptLow func(int64) (err error)
	searchDeptLow = func(dpetId int64) (err error) {
		var deptResult = new(dingOfficialModel.DeptResult) //官方接口的查询结果字段 每一次循环创建一个新的结构
		/*
			调用官方接口比较复杂的一种方法
		*/
		var resp *http.Response
		var body []byte
		URL := "https://oapi.dingtalk.com/topapi/v2/department/listsub?access_token=" + accessToken
		client := &http.Client{Transport: &http.Transport{ //对客户端进行一些配置
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //设置为true , 跳过证书验证
			},
		}, Timeout: time.Duration(time.Second * 9)} //超过9秒取消请求,并且返回一个错误
		deptBody := struct {
			DeptID int64 `json:"dept_id"`
		}{
			DeptID: dpetId,
		}
		//结构体对象序列化
		bodymarshal, err := json.Marshal(&deptBody)
		if err != nil {
			zap.L().Error("json.Marshal(&deptId)")
			return
		}
		reqBody := strings.NewReader(string(bodymarshal))
		//放入具体的request中的
		request, err := http.NewRequest(http.MethodPost, URL, reqBody)
		if err != nil {
			zap.L().Error("http.NewRequest(http.MethodPost, URL, reqBody) is failed", zap.Error(err))
			return
		}
		resp, err = client.Do(request)
		if err != nil {
			zap.L().Error("client.Do(request) is failed", zap.Error(err))
			return
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		//把请求到的结构反序列化到专门接受返回值的对象上面
		err = json.Unmarshal(body, &deptResult)
		if err != nil {
			zap.L().Error("json.Unmarshal(body, &deptResult) is failed", zap.Error(err))
			return
		}
		departments := deptResult.Result //departments查询结果中的部门信息
		deptSearchResult = append(deptSearchResult, departments...)
		if len(departments) > 0 {
			for index, _ := range departments {
				departmentList := make([]deptsModel.TbDept, 0)
				err = searchDeptLow(departments[index].DeptId)
				if err != nil {
					zap.L().Error("searchDeptLow(departments[index].DeptId) is failed", zap.Error(err))
					return err
				}
				deptSearchResult = append(deptSearchResult, departmentList...)
			}
		}
		return nil
	}
	//searchDeptLow查询最低一级的部门 end

	err := searchDeptLow(1) //初始化部门id为1

	//调用service层
	err = dingOfficialService.CreateDeptInformationService(deptSearchResult)
	if err != nil {
		controller.ResponseError(context, controller.CodeServerBusy)
		zap.L().Error("dingOfficialService.CreateDeptInformationService(deptSearchResult) is failed", zap.Error(err))
		return
	}
	controller.ResponseSuccess(context)
}

// GetListUser 得到部门全部用户信息
func (o OfficialController) GetListUser(context *gin.Context) {
	//1.得到部门id 和 accessToken
	accessToken := dingToken.GetOfficialAccessToken()
	if len(accessToken) == 0 {
		context.JSON(http.StatusOK, "权限不足")
		return
	}
	err, deptIds := dingOfficialService.SelectDeptIdsService()
	if err != nil {
		zap.L().Error("dingOfficialService.SelectDeptIdsService() is failed", zap.Error(err))
		return
	}
	//2.首先调用官方接口通过dept_id获取部门用户userId的列表
	var useridListResult []dingOfficialModel.UserIdListResultByDeptId //useridListResult 最终的查询结果
	err, useridListResult = o.DingGetUserIdList(deptIds, accessToken)
	if err != nil {
		zap.L().Error("o.DingGetUserIdList(deptIds) is failed", zap.Error(err))
	}
	//fmt.Println(useridListResult)

	//3.拿到部门用户的userId列表之后 调用官方接口通过部门id查询用户的完整信息
	err, usersResult := o.DingGetUsersInformation(useridListResult, accessToken) // usersResult 从官方接口查询的人员信息的最终结果
	//fmt.Println("#####", usersResult)
	if err != nil {
		zap.L().Error("o.DingGetUsersInformation(useridListResult, accessToken) is failed", zap.Error(err))
		return
	}
	err = dingOfficialService.CreateUserInformationService(usersResult)
	if err != nil {
		zap.L().Error("dingOfficialService.CreateUserInformationService(usersResult) is error", zap.Error(err))
		controller.ResponseError(context, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(context)
}
