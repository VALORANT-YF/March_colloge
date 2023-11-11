package deptsModel

import (
	"gorm.io/gorm"
)

// TbDept 部门表 表名tb_dept
type TbDept struct {
	*gorm.Model
	IsWriteBooks uint8  `json:"is_write_books,omitempty" binding:"oneof=1 0"` //是否需要写简书博客  参数只能为0,1 1表示需要写简书,0表示不需要
	ParentId     int32  `json:"parent_id,omitempty"`                          //父部门id
	DeptId       int64  `json:"dept_id,omitempty"`                            //部门id
	Name         string `json:"name,omitempty"`                               //部门名称
}

type LikeDept struct {
	Name         string `json:"name,omitempty"`
	DeptId       int64  `json:"dept_id,omitempty"`
	IsWriteBooks uint8  `json:"is_write_books,omitempty"`
}

type TbDeptName struct {
	Name string `json:"name"` //部门名称
}

func (t TbDept) TableName() string {
	return "tb_dept"
}

func (l LikeDept) TableName() string {
	return "tb_dept"
}
