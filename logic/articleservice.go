package logic

import (
	"college/dao/mysql"
	"college/models/bookBlogArticle"
	"college/models/usersModel"
	"fmt"
	"strconv"
	"strings"
)

// GetArticleUserInformation 拿到用户信息 包括简书主页链接 , 用户姓名 , 部门名称
func GetArticleUserInformation() (error, []bookBlogArticle.ArticleUserResult) {
	var result []bookBlogArticle.ArticleUserResult
	err, articleUserInformations := mysql.SelectBookAddressPerson() //articleUserInformations 是用户简书主页地址 , 姓名 , 部门编号
	if err != nil {
		return err, result
	}
	//fmt.Println(articleUserInformations)

	for i := 0; i < len(articleUserInformations); i++ {
		//根据部门编号拿到部门名称
		deptList := articleUserInformations[i].DeptIdList
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
			return err, result
		}
		deptNameStr := strings.Join(deptName, " ")
		var bookInformation bookBlogArticle.ArticleUserResult
		bookInformation.Name = articleUserInformations[i].Name
		bookInformation.BookAddress = articleUserInformations[i].BooksAddress
		bookInformation.DeptName = deptNameStr
		bookInformation.Mobile = articleUserInformations[i].Mobile
		bookInformation.BlogAddress = articleUserInformations[i].BlogAddress
		result = append(result, bookInformation)
	}
	return nil, result
}

// InsertAllArticleService 插入所有需要填写简书的用户的简书
func InsertAllArticleService(allArticle []bookBlogArticle.TbBookArticle) error {
	//每周更新一次 , 在插入之前 , 先删除原本的简书数据
	err := mysql.DropArticle()
	if err != nil {
		return err
	}
	return mysql.InsertBookArticle(allArticle)
}

// InsertAllPersonBlogService 插入本周的简书
func InsertAllPersonBlogService(allPersonBlog []bookBlogArticle.TbBlog) error {
	//每周更新一次 , 在插入之前 , 先删除原本的博客数据
	err := mysql.DropBlog()
	if err != nil {
		return err
	}
	return mysql.InsertBlog(allPersonBlog)
}

// SelectAllPersonBookService 查询所有人的简书文章
func SelectAllPersonBookService() (error, []bookBlogArticle.ViewArticle) {
	err, queryResult := mysql.SelectAllPersonBook()
	if err != nil {
		return err, nil
	}

	// 创建一个map集合，合并部门相同的人
	viewMap := make(map[string][]bookBlogArticle.TbBookArticle)
	for _, value := range queryResult {
		if _, ok := viewMap[value.DeptName]; !ok {
			viewMap[value.DeptName] = make([]bookBlogArticle.TbBookArticle, 0)
		}
		viewMap[value.DeptName] = append(viewMap[value.DeptName], value)
	}

	// 将合并后的数据转化为需要的格式
	var result []bookBlogArticle.ViewArticle
	for key, value := range viewMap {
		temp := bookBlogArticle.ViewArticle{
			DeptName:    key,
			TypeArticle: value,
		}
		result = append(result, temp)
	}

	return nil, result // 返回nil作为错误，因为没有出现错误
}

// SelectAllPersonBlogService 查询所有人的博客
func SelectAllPersonBlogService() (error, []bookBlogArticle.ViewBlog) {
	err, queryResult := mysql.SelectAllPersonBlog()
	if err != nil {
		return err, nil
	}
	// 创建一个map集合，合并部门相同的人
	viewMap := make(map[string][]bookBlogArticle.TbBlog)
	for _, value := range queryResult {
		if _, ok := viewMap[value.DeptName]; !ok {
			viewMap[value.DeptName] = make([]bookBlogArticle.TbBlog, 0)
		}
		viewMap[value.DeptName] = append(viewMap[value.DeptName], value)
	}

	// 将合并后的数据转化为需要的格式
	var result []bookBlogArticle.ViewBlog
	for key, value := range viewMap {
		temp := bookBlogArticle.ViewBlog{
			DeptName: key,
			TypeBlog: value,
		}
		result = append(result, temp)
	}

	return nil, result // 返回nil作为错误，因为没有出现错误
}

// ElectExcellentUserService 评选优秀用户的简书或者博客
func ElectExcellentUserService(excellent bookBlogArticle.Excellent) (result bookBlogArticle.ExcellentArticle, err error) {
	//开启事务
	tx := mysql.DB.Begin()
	if tx.Error != nil {
		return result, tx.Error
	}
	//延迟函数用户处理事务提交 或者回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//根据电话号码查找简书主人的信息
	var queryResult usersModel.TbUser
	//评选或取消优秀  1 为简书
	if excellent.BookOrBlog == 1 {
		//根据id 拿到简书主人的信息
		err, userInformation := mysql.SelectBookById(excellent.Id)
		if err != nil {
			return result, err
		}
		//拿到简书主人的电话号码
		mobile := userInformation.Mobile
		err, queryResult = mysql.SelectInformationByTel(mobile)
		//是否评选为优秀简书
		if err = mysql.UpdateExcellentArticle(excellent.Id, excellent.IsTop, tx); err != nil {
			tx.Rollback()
			return result, err
		}
	} else {
		//根据id 拿到简书主人的信息
		err, userInformation := mysql.SelectBlogById(excellent.Id)
		if err != nil {
			return result, err
		}
		//拿到简书主人的电话号码
		mobile := userInformation.Mobile
		err, queryResult = mysql.SelectInformationByTel(mobile)
		//是否评选为优秀博客
		if err = mysql.UpdateExcellentBlog(excellent.Id, excellent.IsTop, tx); err != nil {
			tx.Rollback()
			return result, err
		}
	}
	fmt.Println("人员信息", queryResult)
	if err != nil {
		return result, err
	}
	var newCount uint32
	if excellent.IsTop == 1 {
		newCount = queryResult.ExcellentCount + 1
		fmt.Println("原次数", queryResult.ExcellentCount)
		fmt.Println("优秀简书次数加1", newCount)
	} else {
		newCount = queryResult.ExcellentCount - 1
		fmt.Println("原次数", queryResult.ExcellentCount)
		fmt.Println("优秀简书次数减1", newCount)
	}
	if err = mysql.UpdateUserExcellentCount(queryResult.ID, newCount, tx); err != nil {
		tx.Rollback()
		return result, err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return result, err
	}
	result.IsTop = int(excellent.IsTop)
	result.ID = int(excellent.Id)
	return result, err
}

// IsWriteBookService 部门是否需要写简书
func IsWriteBookService(name string) uint8 {
	return mysql.SelectIsWrite(name)
}

// IsWriteBlogService 部门是否需要写简书
func IsWriteBlogService(name string) uint8 {
	return mysql.SelectIsWriteBlog(name)
}
