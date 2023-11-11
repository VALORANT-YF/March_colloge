package dingOfficialControllers

import (
	"bytes"
	"college/models/dingOfficialModel"
	"fmt"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"net/http"
)

type Ding struct{}

//向官方发的请求放在此处

func (d Ding) DingGetUserIdList(deptIds []int64, accessToken string) (err error, useridListResult []dingOfficialModel.UserIdListResultByDeptId) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/user/listid?access_token=%s", accessToken)
	var response *http.Response
	for _, deptId := range deptIds {
		var userIdResult dingOfficialModel.UserIdListResultByDeptId //userIdResult 单个部门下面的user_id_list 每一次循环时创建一个新的结构
		//封装请求参数
		deptIdParams := struct {
			DeptId int64 `json:"dept_id"`
		}{
			DeptId: deptId,
		}
		//将请求参数转化为JSON字符串
		deptIdParamsJson, err := json.Marshal(deptIdParams)
		reqBody := bytes.NewBuffer([]byte(deptIdParamsJson))
		if err != nil {
			zap.L().Error("json.Marshal(deptIdParams) is failed", zap.Error(err))
			return err, useridListResult
		}
		response, err = http.Post(url, "application/json", reqBody)
		if err = json.NewDecoder(response.Body).Decode(&userIdResult); err != nil {
			zap.L().Error("json.NewDecoder(response.Body).Decode(&useridListResult)", zap.Error(err))
			return err, useridListResult
		}
		//userIdResult.Result.UseridList != 0 的部门下面有用户
		if len(userIdResult.Result.UseridList) != 0 {
			userIdResult.DeptId = deptId
			useridListResult = append(useridListResult, userIdResult)
		}
	}
	defer response.Body.Close()
	return err, useridListResult
}

// DingGetUsersInformation 查询部门用户完成信息
func (d Ding) DingGetUsersInformation(useridListResult []dingOfficialModel.UserIdListResultByDeptId, accessToken string) (err error, usersInformationResult []dingOfficialModel.UsersInformationResult) {
	var response *http.Response
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/v2/user/list?access_token=%s", accessToken)
	for i := 0; i < len(useridListResult); i++ {
		deptId := useridListResult[i].DeptId                          //部门id
		var usersInformation dingOfficialModel.UsersInformationResult // 在每次迭代中创建一个新的结构
		//封装请求参数
		userParams := struct {
			DeptId int64  `json:"dept_id"`
			Cursor uint8  `json:"cursor"`
			Size   uint32 `json:"size"`
		}{
			DeptId: deptId,
			Cursor: 0,
			//len(useridListResult[i].Result.UseridList) 为每一个部门下相应的人数
			Size: uint32(len(useridListResult[i].Result.UseridList)),
		}
		//将请求参数转化为JSON字符串
		userParamsJson, err := json.Marshal(&userParams)
		if err != nil {
			zap.L().Error("json.Marshal(&userParams) is failed", zap.Error(err))
			return err, usersInformationResult
		}
		reqBody := bytes.NewBuffer([]byte(userParamsJson))
		response, err = http.Post(url, "application/json", reqBody)
		if err = json.NewDecoder(response.Body).Decode(&usersInformation); err != nil {
			zap.L().Error("json.NewDecoder(response.Body).Decode(&useridListResult)", zap.Error(err))
			return err, usersInformationResult
		}
		//if deptId == 848512514 { //乐知五期
		//	fmt.Println("####", true)
		//	fmt.Println("@@@", len(useridListResult[i].Result.UseridList))
		//	fmt.Println("!!!!!", usersInformation)
		//}
		usersInformation.DeptId = deptId
		//fmt.Print(usersInformation)
		usersInformationResult = append(usersInformationResult, usersInformation)
	}
	return err, usersInformationResult
}
