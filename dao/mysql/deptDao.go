package mysql

import "college/models/deptsModel"

// UpdateDeptIsWrite 修改部门是否需要写简书的状态
func UpdateDeptIsWrite(isWrite uint8, deptId int64) error {
	return DB.Where("dept_id = ?", deptId).Table("tb_dept").Update("is_write_books", isWrite).Error
}

// SelectAllDept 查询所有子部门名称和部门id以及父部门id
func SelectAllDept() (error, []deptsModel.TbDept) {
	var lowDeptInformation []deptsModel.TbDept
	//先查询子部门的部门名称和部门id
	err := DB.Select("parent_id , dept_id , name , is_write_books").Find(&lowDeptInformation).Error
	return err, lowDeptInformation
}

// SelectDeptNameByParentId 根据parent_id 拿到父部门名称
func SelectDeptNameByParentId(parentId int32) (error, deptsModel.TbDept) {
	var parentDeptName deptsModel.TbDept
	err := DB.Where("dept_id = ?", parentId).Select("name").Find(&parentDeptName).Error
	return err, parentDeptName
}

// SelectNeedWriteDept 查询出需要写简书的部门
func SelectNeedWriteDept() (err error, tbDepts []deptsModel.TbDept) {
	err = DB.Where("is_write_books = 1").Find(&tbDepts).Error
	return
}

// SelectIsWrite 根据部门名称查询部门是否需要写简书
func SelectIsWrite(name string) uint8 {
	var isWirte uint8
	DB.Select("is_write_books").Where("name = ?", name).Table("tb_dept").Find(&isWirte)
	return isWirte
}

// SelectIsWriteBlog 查询部门是否需要写博客
func SelectIsWriteBlog(name string) uint8 {
	var isWirte uint8
	DB.Select("is_write_books").Where("name = ?", name).Table("tb_dept").Find(&isWirte)
	return isWirte
}
