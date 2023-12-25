package timeTask

import (
	"college/controller"
	"college/controller/dingOfficialControllers"
	"college/logic"
	"college/models/bookBlogArticle"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

type Test struct {
}

// 测试使用
func (t Test) GetUserListTest(context *gin.Context) {
	var message strings.Builder
	zap.L().Info("开始获取本周简书未写名单和优秀名单")
	err, noWriteList := logic.SelectNoWriteUserService() //简书博客未写的人员名单
	if err != nil {
		zap.L().Error("logic.SelectNoWriteUserService() is error", zap.Error(err))
	}
	message.WriteString("\n[吃瓜]简书博客未写人员名单:\n")
	for _, value := range noWriteList {
		message.WriteString(value.OneDeptName + ":")
		message.WriteString("\n")
		for _, data := range value.NoWriteUserList {
			message.WriteString(data.Name + " ")
		}
		message.WriteString("\n")
	}
	message.WriteString("\n")

	excellentList := logic.SelectExcellentPersonService() //优秀简书的名单

	message.WriteString("\n[元气满满]优秀简书博客名单:\n")
	for _, value := range excellentList {
		message.WriteString(value.DeptName + ":")
		message.WriteString("\n")
		for _, data := range value.TypeName {
			//类型断言
			switch v := data.(type) {
			case bookBlogArticle.TbBookArticle:
				message.WriteString(v.Name + "\t" + "简书: " + v.BookTitle + "\n" + v.ArticleUrl + "\n")
			case bookBlogArticle.TbBlog:
				message.WriteString(v.Name + "\t" + "博客" + v.BookTitle + "\n" + v.BlogUrl + "\n")
			default:
				fmt.Println("Unknown Type")
			}
		}
	}
	message.WriteString("\n")
	//排行榜
	err, resultExcellent := logic.SelectExcellentCountService()
	if err != nil {
		zap.L().Error("logic.SelectExcellentCountService()", zap.Error(err))
		return
	}

	err, resultNoWriteCount := logic.SelectNoWriteCount()
	if err != nil {
		zap.L().Error("logic.SelectNoWriteUserService()", zap.Error(err))
		return
	}

	message.WriteString("[小蜜蜂]优秀简书次数排行榜(前5名):\n")
	for _, data := range resultExcellent {
		message.WriteString("姓名: " + data.Name + "\t次数: " + fmt.Sprintf("%v", data.ExcellentCount))
		message.WriteString("\n")
	}

	message.WriteString("\n[疑问]未写简书次数排行榜(前3名):\n")
	for _, data := range resultNoWriteCount {
		message.WriteString("姓名: " + data.Name + "\t次数: " + fmt.Sprintf("%v", data.NotWrittenCount))
		message.WriteString("\n")
	}
	message.WriteString("\n[撒花]注意： \n\n[三多]未填写简书者为组内贡献班费用(6元[奶茶]），并且如果在周二24：00前没有对简书进行补录，那么将罚(12元[奶茶][奶茶]）。费用将由财务部收取。\n\n[三多]提交简书链接时务必检查简书链接是否正常，不允许出现\"页面不存在\"或者\"正在审核\"的情况，如果出现，一律按照未填写处理 \n\n[三多]如有特殊情况可找纪检部:\n李壮 电话:19838787058\n\n纪检部通知[广播][广播]断水断电不断简书，为谋为事必为总结[送花花][跳舞]。")
	//调用机器人发送消息 拿到需要发送消息的机器人的token
	err, robotToken := logic.GetRobotTokenList()
	if err != nil {
		zap.L().Error("获取机器人token失败")
		return
	}
	for _, data := range robotToken {
		err = dingOfficialControllers.SendTextMessage(message.String(), data.WebhookURL)
		if err != nil {
			zap.L().Error("机器人发送简书博客消息失败")
			return
		}
	}
	zap.L().Info("结束获取本周简书未写名单和优秀名单")
}

func (t Test) RemindTest(context *gin.Context) {
	var message strings.Builder
	zap.L().Info("机器人提醒任务")
	err, noWriteList := logic.SelectNoWriteUserService() //简书博客未写的人员名单
	if err != nil {
		zap.L().Error("logic.SelectNoWriteUserService() is error", zap.Error(err))
	}
	message.WriteString("纪检部通知[广播][广播]:  \n断水断电不断简书，为谋其事必为总结[爱意]  \n[钉子]各期负责人于下周一晚上20：00前在钉钉简书小程序中标记优秀简书\n  \n[钉子]检查为机器人检查，大家要及时发表文章\n\n  [忍者][忍者]注意\n[对勾]简书严禁抄袭，坚持原创，我们会根据相关字段进行严格检查的哦[爱意]\n[对勾]纪检部同时也会对简书进行抽查[猫咪]\n[对勾]简书字数不能低于400字\n[钉子][钉子]重点！！！到周日20：00后财务部人员会在大群里面对简书未完成人员发起群收款[惊愕][惊愕]\n大家要注意了哦！！\n\n[灵感][灵感]提醒:\n [对勾]简书以及博客的时间为本周内，否则会被标记为未登记！\n\n[爱意]希望大家多多参与简书投稿及评论互动[捧脸]并在此相互学习和借鉴哟[猫咪][猫咪]\n")
	message.WriteString("\n[吃瓜]简书博客未写人员名单:\n")
	for _, value := range noWriteList {
		message.WriteString(value.OneDeptName + ":")
		message.WriteString("\n")
		for _, data := range value.NoWriteUserList {
			message.WriteString(data.Name + "\t")
		}
		message.WriteString("\n")
	}
	message.WriteString("\n")
	//调用机器人发送消息 拿到需要发送消息的机器人的token
	err, robotToken := logic.GetRobotTokenList()
	if err != nil {
		zap.L().Error("获取机器人token失败")
		return
	}
	for _, data := range robotToken {
		err = dingOfficialControllers.SendTextMessage(message.String(), data.WebhookURL)
		if err != nil {
			zap.L().Error("机器人发送提醒消息失败")
			return
		}
	}
	zap.L().Info("机器人发送提醒消息成功")
}

var articleController controller.ArticleController
var wg sync.WaitGroup

func TimeTask(context *gin.Context) error {
	crontab := cron.New(cron.WithSeconds()) //时间精确到秒
	//定义定时器调用的任务函数

	//设置时间
	//specTest1 := "0 42 14 ? * 6"
	//specTest2 := "0 45 16 ? * 6"
	specArticle := "0 30 23 * * 0" //每周周日 11点半 一次

	specMessage := "0 30 17 ? * 2" //周二下午上17点半

	//两个提醒消息 提醒之前需要先爬取简书博客
	specFriday := "0 30 15 * * 5" //周五下午15:30
	specSunday := "0 30 16 * * 0" //周日下午16:30
	//添加定时任务
	_, err := crontab.AddFunc(specArticle, func() {
		//爬取简书博客的定时任务 正式爬取
		articleBlog(context)
	})

	if err != nil {
		return err
	}
	//添加定时任务
	_, err = crontab.AddFunc(specMessage, func() {
		//获取简书未写名单以及优秀人员名单的任务
		GetUserList()
	})

	//提醒任务,两个时间段
	_, err = crontab.AddFunc(specFriday, func() {
		//周五下午3点半一次 只是提醒
		RemindMessage()
	})

	if err != nil {
		zap.L().Error("周五下午的提醒定时任务启动失败")
		return err
	}
	_, err = crontab.AddFunc(specSunday, func() {
		articleBlog(context)
		//周日下午3:30提醒一次 提醒之前先爬取 发送未写此时未写简书的名单
		RemindMessage()
	})
	if err != nil {
		zap.L().Error("周日下午的提醒定时任务启动失败")
		return err
	}

	//启动定时器
	crontab.Start()
	return nil
}

// articleBlog 定时任务,每周周日晚上11:30.调用爬取简书博客的接口
func articleBlog(context *gin.Context) {
	zap.L().Info("简书博客爬取开始")
	wg.Add(2)
	go taskBlog(context)
	go taskArticle(context)
	wg.Wait()
	zap.L().Info("简书博客爬取成功")
}

func taskArticle(context *gin.Context) {
	articleController.UserArticle(context)
	zap.L().Info("爬取简书成功")
	wg.Done()
}

func taskBlog(context *gin.Context) {
	articleController.UserBlog(context)
	zap.L().Info("爬取博客成功")
	wg.Done()
}

// 机器人发送消息

// GetUserList 简书未写名单和优秀简书名单 每周周二下午17:30
func GetUserList() {
	var message strings.Builder
	zap.L().Info("开始获取本周简书未写名单和优秀名单")
	err, noWriteList := logic.SelectNoWriteUserService() //简书博客未写的人员名单
	if err != nil {
		zap.L().Error("logic.SelectNoWriteUserService() is error", zap.Error(err))
	}
	message.WriteString("\n[吃瓜]简书博客未写人员名单:\n")
	for _, value := range noWriteList {
		message.WriteString(value.OneDeptName + ":")
		message.WriteString("\n")
		for _, data := range value.NoWriteUserList {
			message.WriteString(data.Name + " ")
		}
		message.WriteString("\n")
	}
	message.WriteString("\n")

	excellentList := logic.SelectExcellentPersonService() //优秀简书的名单

	message.WriteString("\n[元气满满]优秀简书博客名单:\n")
	for _, value := range excellentList {
		message.WriteString(value.DeptName + ":")
		message.WriteString("\n")
		for _, data := range value.TypeName {
			//类型断言
			switch v := data.(type) {
			case bookBlogArticle.TbBookArticle:
				message.WriteString(v.Name + "\t" + "简书: " + v.BookTitle + "\n" + v.ArticleUrl + "\n")
			case bookBlogArticle.TbBlog:
				message.WriteString(v.Name + "\t" + "博客" + v.BookTitle + "\n" + v.BlogUrl + "\n")
			default:
				fmt.Println("Unknown Type")
			}
		}
	}
	message.WriteString("\n")
	//排行榜
	err, resultExcellent := logic.SelectExcellentCountService()
	if err != nil {
		zap.L().Error("logic.SelectExcellentCountService()", zap.Error(err))
		return
	}

	err, resultNoWriteCount := logic.SelectNoWriteCount()
	if err != nil {
		zap.L().Error("logic.SelectNoWriteUserService()", zap.Error(err))
		return
	}

	message.WriteString("[小蜜蜂]优秀简书次数排行榜(前5名):\n")
	for _, data := range resultExcellent {
		message.WriteString("姓名: " + data.Name + "\t次数: " + fmt.Sprintf("%v", data.ExcellentCount))
		message.WriteString("\n")
	}

	message.WriteString("\n[疑问]未写简书次数排行榜(前3名):\n")
	for _, data := range resultNoWriteCount {
		message.WriteString("姓名: " + data.Name + "\t次数: " + fmt.Sprintf("%v", data.NotWrittenCount))
		message.WriteString("\n")
	}
	message.WriteString("\n[撒花]注意： \n\n[三多]未填写简书者为组内贡献班费用(6元[奶茶]），并且如果在周二24：00前没有对简书进行补录，那么将罚(12元[奶茶][奶茶]）。费用将由财务部收取。\n\n[三多]提交简书链接时务必检查简书链接是否正常，不允许出现\"页面不存在\"或者\"正在审核\"的情况，如果出现，一律按照未填写处理 \n\n[三多]如有特殊情况可找纪检部:\n李壮 电话:19838787058\n\n纪检部通知[广播][广播]断水断电不断简书，为谋为事必为总结[送花花][跳舞]。")
	//调用机器人发送消息 拿到需要发送消息的机器人的token
	err, robotToken := logic.GetRobotTokenList()
	if err != nil {
		zap.L().Error("获取机器人token失败")
		return
	}
	for _, data := range robotToken {
		err = dingOfficialControllers.SendTextMessage(message.String(), data.WebhookURL)
		if err != nil {
			zap.L().Error("机器人发送简书博客消息失败", zap.Error(err))
			continue
		}
	}
	zap.L().Info("结束获取本周简书未写名单和优秀名单")
}

// RemindMessage 机器人每周周五和周日提醒写简书博客
func RemindMessage() {
	var message strings.Builder
	zap.L().Info("机器人提醒任务")
	err, noWriteList := logic.SelectNoWriteUserService() //简书博客未写的人员名单
	if err != nil {
		zap.L().Error("logic.SelectNoWriteUserService() is error", zap.Error(err))
	}
	message.WriteString("纪检部通知[广播][广播]:  \n断水断电不断简书，为谋其事必为总结[爱意]  \n[钉子]各期负责人于下周一晚上20：00前在钉钉简书小程序中标记优秀简书\n  \n[钉子]检查为机器人检查，大家要及时发表文章\n\n  [忍者][忍者]注意\n[对勾]简书严禁抄袭，坚持原创，我们会根据相关字段进行严格检查的哦[爱意]\n[对勾]纪检部同时也会对简书进行抽查[猫咪]\n[对勾]简书字数不能低于400字\n[钉子][钉子]重点！！！到周日20：00后财务部人员会在大群里面对简书未完成人员发起群收款[惊愕][惊愕]\n大家要注意了哦！！\n\n[灵感][灵感]提醒:\n [对勾]简书以及博客的时间为本周内，否则会被标记为未登记！\n\n[爱意]希望大家多多参与简书投稿及评论互动[捧脸]并在此相互学习和借鉴哟[猫咪][猫咪]\n")
	//如果现在时间是周五 不去拼接未写人员名单
	currentTime := time.Now() //获取现在时间
	//判断当前时间是周几
	dayOfWeek := currentTime.Weekday()
	//如果当前时间不是周五 , 而是周日拼接简书博客未写名单
	if dayOfWeek != time.Friday {
		message.WriteString("\n[吃瓜]简书博客未写人员名单:\n")
		for _, value := range noWriteList {
			message.WriteString(value.OneDeptName + ":")
			message.WriteString("\n")
			for _, data := range value.NoWriteUserList {
				message.WriteString(data.Name + "\t")
			}
			message.WriteString("\n")
		}
	}
	message.WriteString("\n")
	//调用机器人发送消息 拿到需要发送消息的机器人的token
	err, robotToken := logic.GetRobotTokenList()
	if err != nil {
		zap.L().Error("获取机器人token失败")
		return
	}
	for _, data := range robotToken {
		err = dingOfficialControllers.SendTextMessage(message.String(), data.WebhookURL)
		if err != nil {
			zap.L().Error("机器人发送提醒消息失败", zap.Error(err))
			continue
		}
	}
	zap.L().Info("机器人发送提醒消息成功")
}
