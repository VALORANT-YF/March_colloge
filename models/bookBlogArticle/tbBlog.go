package bookBlogArticle

import "gorm.io/gorm"

// TbBlog 用户博客文章表 只存一周的文章
type TbBlog struct {
	gorm.Model
	BookTitle   string `json:"book_title,omitempty"`                    //博客标题
	BookArticle string `json:"book_article,omitempty" gorm:"type:text"` //使用长文本类型存储博客文章
	DeptName    string `json:"dept_name,omitempty"`                     //博客主人所属的部门名称
	Name        string `json:"name,omitempty"`                          //博客文章的主人的姓名
	Mobile      string `json:"mobile,omitempty"`                        //博客文章主任的电话
	IsTop       string `json:"is_top"`                                  //博客是否置顶
	BlogUrl     string `json:"blog_url,omitempty"`                      //原文链接
}

func (t TbBlog) TableName() string {
	return "tb_blog"
}
