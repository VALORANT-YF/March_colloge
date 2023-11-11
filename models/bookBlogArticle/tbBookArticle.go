package bookBlogArticle

import "gorm.io/gorm"

// TbBookArticle 用户简书文章表 只存一周的文章
type TbBookArticle struct {
	gorm.Model
	BookTitle   string `json:"book_title,omitempty"`                    //简书标题
	BookArticle string `json:"book_article,omitempty" gorm:"type:text"` //使用长文本类型存储简书文章
	DeptName    string `json:"dept_name,omitempty"`                     //简书博客主任所属的部门名称
	Name        string `json:"name,omitempty"`                          //简书博客文章的主人的姓名
	Mobile      string `json:"mobile,omitempty"`                        //简书博客文章主任的电话
	IsTop       string `json:"is_top"`                                  //文章是否是优秀简书(优秀简书置顶)
	ArticleUrl  string `json:"article_url,omitempty"`                   //原文链接
}

func (t TbBookArticle) TableName() string {
	return "tb_book_article"
}
