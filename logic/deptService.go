package logic

import (
	"college/dao/mysql"
	"college/models/deptsModel"
	"strconv"
	"strings"
)

// SelectDeptPersonInformationService 查询出正确的部门结构和部门人员,返回
/*
 如:十一期强化班:
	1组: 1号
*/
func SelectDeptPersonInformationService(name string) (error, []deptsModel.ResultHigh) {
	//创建map 集合 , 将HighDept相同的合并
	resultMap := make(map[string][]deptsModel.Result)

	var resultView []deptsModel.DeptPersonInformation //最终结果
	err, allDeptInformation := mysql.SelectAllDept()
	if err != nil {
		return err, nil
	}
	var deptPersonInformation []deptsModel.DeptPersonInformation
	//根据父部门拿到部门关系
	for _, data := range allDeptInformation {
		var deptInformaiton deptsModel.DeptPersonInformation
		//fmt.Println("parent_id", data.ParentId)
		//根据父部门id拿到父部门的名称
		err, parentDeptName := mysql.SelectDeptNameByParentId(data.ParentId)
		if err != nil {
			return err, nil
		}
		deptInformaiton.HighDept = parentDeptName.Name   //高级别的部门名称
		deptInformaiton.DeptId = data.DeptId             //部门id
		deptInformaiton.DeptName = data.Name             //子部门名称
		deptInformaiton.IsWriteBooks = data.IsWriteBooks //部门是否需要写简书
		deptPersonInformation = append(deptPersonInformation, deptInformaiton)
	}
	//拿到人员信息
	err, allPersonInformation := mysql.SelectAllPersonInformation(name)
	for _, dept := range deptPersonInformation {
		for _, data := range allPersonInformation {
			var personInformation deptsModel.PersonInformation
			deptStr := strconv.FormatInt(dept.DeptId, 10)
			//如果用户部门列表中包含此部门,将用户信息插入到该部门中
			if strings.Contains(data.DeptIdList, deptStr) {
				//如果部门下面没有人,直接跳过本次循环
				personInformation.UserId = data.Userid        //用户id
				personInformation.Name = data.Name            //用户名
				personInformation.Mobile = data.Mobile        //用户电话
				personInformation.IsBoss = data.IsBoss        //是否是管理员
				personInformation.BlogUrl = data.BlogAddress  //博客主页地址
				personInformation.BookUrl = data.BooksAddress //简书主页地址
				dept.TypeName = append(dept.TypeName, personInformation)
			}
		}
		resultView = append(resultView, dept)
	}

	for _, data := range resultView {
		//如果这个部门下面没有人员 不去渲染
		if data.TypeName == nil {
			continue
		}
		//该部门就是最高部门
		if data.HighDept == "" {
			data.HighDept = data.DeptName
		}
		//如果部门存在
		resultMap[data.HighDept] = append(resultMap[data.HighDept], deptsModel.Result{
			DeptName:     data.DeptName,     //部门名称
			DeptId:       data.DeptId,       //部门id
			IsWriteBooks: data.IsWriteBooks, //部门是否需要写简书
			TypeName:     data.TypeName,
		})
	}

	var endResults []deptsModel.ResultHigh

	for highDeptName, data := range resultMap {
		var endResult deptsModel.ResultHigh
		endResult.HighDeptName = highDeptName
		endResult.TypeNameEnd = data
		endResults = append(endResults, endResult)
	}

	return err, endResults
}
