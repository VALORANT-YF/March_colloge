package adminGroups

import (
	"college/controller"
	middles "college/middlewares"
	"college/pkg/timeTask"
	"github.com/gin-gonic/gin"
)

var adminController controller.AdminController
var test timeTask.Test

// AdminRouters 管理员路由
func AdminRouters(r *gin.Engine) {

	//superAdmin 超级管理员
	superAdmin := r.Group("/superAdmin")
	{
		superAdmin.POST("/superLogin", adminController.GetTokenController)                   //超级管理员登录
		superAdmin.POST("/allDept", adminController.SelectDeptPersonInformation)             //渲染页面
		superAdmin.GET("/isWriteDept/:dept_id/:is_write_books", adminController.DeptIsWrite) //设置需要写简书的部门
		superAdmin.GET("/updateAdmin/:userid/:is_boss", adminController.UpdateAdmin)         //修改管理员
	}

	//普通管理员
	admin := r.Group("/admin", middles.JWTAuthMiddleware())
	{
		admin.POST("/elect", adminController.ExcellentBookBlog) //普通管理员评选优秀简书
		admin.GET("/noWriteList", adminController.NoWriteUsers) //管理员查看本周简书未写名单
		admin.GET("/lookRobot", adminController.LookRobot)      //管理员查看机器人
		admin.POST("/addRobot", adminController.AddRobot)       //增加机器人
		admin.POST("/updateRobot", adminController.UpdateRobot) //修改机器人
		admin.POST("/dropRobot", adminController.DeleteRobot)   //删除机器人
	}

	//管理相关的定时任务
	adminTime := r.Group("/admin-time")
	{
		//下面4个是定时任务 , 不需要前端调用 , 写接口是为了测试
		adminTime.GET("/noWriteList", adminController.NoWriteUsers)          //简书未写名单
		adminTime.GET("/excellent", adminController.ExcellentPersonList)     //本周优秀简书名单
		adminTime.GET("/excellentCount", adminController.ExcellentCountFive) //优秀次数前5
		adminTime.GET("/noWriteCount", adminController.NoWriteCountThree)    //未写次数前3
		adminTime.GET("/test", test.GetUserListTest)                         //测试定时任务 简书未写名单
		adminTime.GET("/test2", test.RemindTest)                             //提醒消息 测试定时任务
	}
}
