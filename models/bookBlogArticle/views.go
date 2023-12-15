package bookBlogArticle

type TypeArticle []TbBookArticle

type TypeBlog []TbBlog

type ViewArticle struct {
	DeptName string `json:"dept_name"`
	TypeArticle
}

type ViewBlog struct {
	DeptName string `json:"dept_name"`
	TypeBlog
}

type TypeName []interface{}

type ViewResult struct {
	DeptName string `json:"dept_name"`
	TypeName TypeName
}

type ExcellentArticle struct {
	ID    int `json:"id"`
	IsTop int `json:"is_top"`
}
