package deptsModel

//需要展示在页面上的部门人员信息架构

type TypeName []PersonInformation

type DeptPersonInformation struct {
	HighDept     string `json:"high_dept,omitempty"` //最高级别的部门名称
	DeptName     string `json:"dept_name,omitempty"` //子部门名称
	DeptId       int64  `json:"dept_id,omitempty"`   //部门id
	IsWriteBooks uint8  `json:"is_write_books"`      //是否需要写简书
	TypeName
}

type PersonInformation struct {
	Name    string `json:"name,omitempty"`     //人员名称
	Mobile  string `json:"mobile,omitempty"`   //人员电话
	IsBoss  bool   `json:"is_boss"`            //是否是管理员
	UserId  string `json:"user_id,omitempty"`  //user_id
	BookUrl string `json:"book_url,omitempty"` //简书主页地址
	BlogUrl string `json:"blog_url,omitempty"` //博客主页地址
}

type Result struct {
	DeptName     string `json:"dept_name,omitempty"` //子部门名称
	DeptId       int64  `json:"dept_id,omitempty"`   //部门id
	IsWriteBooks uint8  `json:"is_write_books"`      //是否需要写简书
	TypeName
}

type TypeNameEnd []Result

type ResultHigh struct {
	HighDeptName string `json:"high_dept_name"`
	TypeNameEnd
}
