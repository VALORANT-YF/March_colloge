package usersModel

//关于用户请求相关的参数封装

// UserAccount 封装用户登录时的请求参数
type UserAccount struct {
	Name     string `json:"name,omitempty" binding:"required"`     //用户名
	Mobile   string `json:"mobile,omitempty" binding:"required"`   //电话号码
	Password string `json:"password,omitempty" binding:"required"` //密码
}

// BlogBookAddress 接收用户插入简书博客主页地址时传入的参数
type BlogBookAddress struct {
	BooksAddress string `json:"books_address,omitempty" binging:"required"` //简书主页地址
	BlogAddress  string `json:"blog_address,omitempty" binding:"required"`  //博客主页地址
}

// UpdateAdminUri 封装修改管理员的查询参数
type UpdateAdminUri struct {
	Userid string `uri:"userid" binding:"required"`
	IsBoss int    `uri:"is_boss"`
}

// UserLookBookOrBlog 用户查询其他人简书或者博客的参数封装
type UserLookBookOrBlog struct {
	IsBookOrBlog int `binding:"oneof=0 1"` //0表示用户点击的文章是简书 , 1表示博客
	Id           int `binding:"required"`  //简书或者博客id
}

// UserLookBookBlogDetails 用户查看的详细内容
type UserLookBookBlogDetails struct {
	Id          int64  `json:"id,omitempty"`           //文章id
	BookTitle   string `json:"book_title,omitempty"`   //文章标题
	BookArticle string `json:"book_article,omitempty"` //文章内容
	DeptName    string `json:"dept_name,omitempty"`    //文章主人所在的部门
	Name        string `json:"name,omitempty"`         //姓名
	Url         string `json:"url,omitempty"`          //原文链接
}

// UserSelfInformation 封装用户个人信息
type UserSelfInformation struct {
	Name         string `json:"name,omitempty"`          //用户姓名
	Avatar       string `json:"avatar,omitempty"`        //头像地址
	DeptName     string `json:"dept_name,omitempty"`     //部门名称
	Mobile       string `json:"mobile,omitempty"`        //电话号码
	BooksAddress string `json:"books_address,omitempty"` //简书主页地址
	BlogAddress  string `json:"blog_address,omitempty"`  //博客主页地址
	Password     string `json:"password,omitempty"`      //用户密码
}

// UserUpdateSelfInformation 封装用户修改个人信息时接收的参数
type UserUpdateSelfInformation struct {
	BooksAddress string `json:"books_address,omitempty"` //简书主页地址
	BlogAddress  string `json:"blog_address,omitempty"`  //博客主页地址
	Password     string `json:"password,omitempty"`      //密码
}
