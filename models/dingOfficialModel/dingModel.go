package dingOfficialModel

import (
	"college/models/deptsModel"
)

// ErrorResult 官方返回的错误结果
type ErrorResult struct {
	Errcode int    `json:"errcode"` //错误状态码
	Errmsg  string `json:"errmsg"`  //错误消息
}

// DeptResult 封装查询部门列表的最后结果
type DeptResult struct {
	ErrorResult
	Result []deptsModel.TbDept `json:"result"`
}

// UserIdListResult 封装查询出来的UserIdList
type UserIdListResult struct {
	UseridList []string `json:"userid_list"`
}

// UserIdListResultByDeptId 封装查询部门用户user_id 的查询结果
type UserIdListResultByDeptId struct {
	ErrorResult
	DeptId int64            `json:"dept_id"`
	Result UserIdListResult `json:"result"`
}

// UserInformation 查询出的个人信息
type UserInformation struct {
	Boss       bool    `json:"boss,omitempty"`         //是否是boss
	Unionid    string  `json:"unionid,omitempty"`      //用户的唯一id
	Mobile     string  `json:"mobile,omitempty"`       //电话号码
	Avatar     string  `json:"avatar,omitempty"`       //头像路径
	Userid     string  `json:"userid,omitempty"`       //用户id
	Name       string  `json:"name,omitempty"`         //用户名
	DeptIdList []int64 `json:"dept_id_list,omitempty"` //用户部门列表集合
}

// UsersInformation 查询出全部人员的个人信息
type UsersInformation struct {
	List []UserInformation `json:"list,omitempty"`
}

// UsersInformationResult 最终查询结果
type UsersInformationResult struct {
	ErrorResult
	DeptId int64            `json:"dept_id,omitempty"`
	Result UsersInformation `json:"result,omitempty"`
}
