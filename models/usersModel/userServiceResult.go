package usersModel

//Service层返回的结果

// NoWriteUser 简书未写的人
type NoWriteUser struct {
	NotWrittenCount uint32 `json:"not_written_count,omitempty"` //简书未写次数
	Mobile          string `json:"mobile,omitempty"`            //电话
	UserId          string `json:"user_id,omitempty"`           //用户id
	Avatar          string `json:"avatar,omitempty"`            //头像
	Name            string `json:"name,omitempty"`              //姓名
	BooksAddress    string `json:"books_address,omitempty"`     //简书地址
	BlogAddress     string `json:"blog_address,omitempty"`      //博客地址
	DeptName        string `json:"dept_name,omitempty"`         //此人所在的部门名称
}

type NoWriteUserList []NoWriteUser

type NoWriteView struct {
	OneDeptName     string `json:"one_dept_name"`
	NoWriteUserList NoWriteUserList
}
