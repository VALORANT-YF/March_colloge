package usersModel

import "gorm.io/gorm"

type TbUser struct {
	gorm.Model
	IsBoss          bool   `json:"is_boss,omitempty"`                   //是否是领导
	ExcellentCount  uint32 `json:"excellent_count,omitempty"`           //优秀简书或者博客的次数
	NotWrittenCount uint32 `json:"not_written_count,omitempty"`         //简书未写次数
	Unionid         string `json:"unionid,omitempty"`                   //用户的唯一id
	Mobile          string `json:"mobile,omitempty" binding:"required"` //电话号码
	Userid          string `json:"userid,omitempty"`                    //用户id
	Avatar          string `json:"avatar,omitempty"`                    //钉钉头像地址
	Name            string `json:"name,omitempty"`                      //姓名
	BooksAddress    string `json:"books_address,omitempty"`             //简书主页地址
	BlogAddress     string `json:"blog_address,omitempty"`              //博客数据地址
	DeptIdList      string `json:"dept_id_List,omitempty"`              //用户所属的部门id 一个用户可以在多个部门下 插入数据库时需要修改为字符串类型
	Password        string `json:"password,omitempty"`                  //密码
}

// TbUserBookBlog 用户简书博客主页地址
type TbUserBookBlog struct {
	BookAddress string `json:"book_address,omitempty"`
	BlogAddress string `json:"blog_address,omitempty"`
}

// TbUserDeptAndName 用户姓名和部门列表
type TbUserDeptAndName struct {
	Name       string `json:"name,omitempty"`         //用户姓名
	DeptIdList string `json:"dept_id_list,omitempty"` //用户的部门列表
}

// AddressPerson 用户简书博客主页地址,姓名,电话,部门编号
type AddressPerson struct {
	BooksAddress string `json:"books_address,omitempty"` //简书主页地址
	BlogAddress  string `json:"blog_address,omitempty"`  //博客主页地址
	Name         string `json:"name,omitempty"`          //姓名
	Mobile       string `json:"mobile,omitempty"`        //电话
	DeptIdList   string `json:"dept_id_list,omitempty"`  //部门编号
}

func (t TbUser) TableName() string {
	return "tb_user"
}

func (a AddressPerson) TableName() string {
	return "tb_user"
}
