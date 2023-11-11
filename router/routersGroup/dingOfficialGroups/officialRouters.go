package dingOfficialGroups

import (
	"college/controller/dingOfficialControllers"
	"github.com/gin-gonic/gin"
)

var officialController dingOfficialControllers.OfficialController

// DingOfficialRouters 钉钉官方接口 , 获取部门人员的全部信息
func DingOfficialRouters(r *gin.Engine) {
	official := r.Group("/ding")
	{
		official.POST("/getDeptList", officialController.GetDeptList) //官方接口,获取架构中的部门列表
		official.POST("/getUsers", officialController.GetListUser)    //官方接口,获取架构中的全部用户的信息
		official.POST("/updateList")                                  //更新部门信息
	}
}
