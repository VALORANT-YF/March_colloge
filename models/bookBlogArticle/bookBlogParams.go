package bookBlogArticle

type ArticleUserResult struct {
	Name        string `json:"name,omitempty"`         //用户姓名
	BookAddress string `json:"book_address,omitempty"` //简书主页地址
	BlogAddress string `json:"blog_address,omitempty"` //博客主页地址
	Mobile      string `json:"mobile,omitempty"`       //用户电话号码
	DeptName    string `json:"dept_name,omitempty"`    //用户所属部门
}

type Excellent struct {
	BookOrBlog uint8 `json:"book_or_blog,omitempty"`          //评选的文章是简书还是博客 1为简书 0 为博客
	Id         int64 `json:"id,omitempty" binding:"required"` //简书或者博客的id
	IsTop      uint8 `json:"is_top,omitempty"`                //简书或者博客是否置顶 1为优秀简书 或者 优秀博客
}
