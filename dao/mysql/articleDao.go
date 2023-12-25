package mysql

import (
	"college/models/bookBlogArticle"
	"college/models/usersModel"
	"fmt"
	"gorm.io/gorm"
)

// 定义表名常量
const (
	tableNameBlog = "tb_blog"
	tableNameBook = "tb_book_article"
)

// InsertBookArticle 插入简书文章
func InsertBookArticle(allPersonArticle []bookBlogArticle.TbBookArticle) error {
	//开启事务
	tx := DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	//遍历数据 执行插入操作
	fmt.Println("%%%%%", len(allPersonArticle))
	for _, data := range allPersonArticle {
		if err := tx.Create(&data).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// DropArticle 清空简书表中的文章
func DropArticle() error {
	return DB.Exec("TRUNCATE TABLE tb_book_article").Error
}

// DropBlog 清空博客表中的文章
func DropBlog() error {
	return DB.Exec("TRUNCATE TABLE tb_blog").Error
}

// InsertBlog 插入每一个人的博客
func InsertBlog(allPersonBlog []bookBlogArticle.TbBlog) error {
	tx := DB.Begin() //开启事务

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	//遍历数据 执行插入操作
	for _, data := range allPersonBlog {
		if err := tx.Create(&data).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// SelectAllPersonBook 查询所有人的简书
func SelectAllPersonBook() (error, []bookBlogArticle.TbBookArticle) {
	var allBooks []bookBlogArticle.TbBookArticle
	err := DB.Select("id,book_title , book_article , dept_name , name , mobile , article_url , is_top").Find(&allBooks).Error
	return err, allBooks
}

// SelectAllPersonBlog 查询所有人的博客
func SelectAllPersonBlog() (error, []bookBlogArticle.TbBlog) {
	var allBlogs []bookBlogArticle.TbBlog
	err := DB.Select("id,book_title , book_article , dept_name , name , mobile , blog_url , is_top").Find(&allBlogs).Error
	return err, allBlogs
}

// SelectBookById 根据id去查询简书的详细内容
func SelectBookById(id int64) (error, bookBlogArticle.TbBookArticle) {
	var oneBookContent bookBlogArticle.TbBookArticle
	err := DB.Select("id,book_title , book_article , dept_name , name , mobile , article_url , is_top").Where("id = ?", id).Find(&oneBookContent).Error
	return err, oneBookContent
}

// SelectBlogById 根据id去查询博客的详细内容
func SelectBlogById(id int64) (error, bookBlogArticle.TbBlog) {
	var oneBlogContent bookBlogArticle.TbBlog
	err := DB.Select("id,book_title , book_article , dept_name , name , mobile , blog_url , is_top").Where("id = ?", id).Find(&oneBlogContent).Error
	return err, oneBlogContent
}

// UpdateExcellentArticle 管理员设置优秀简书 , 或者取消优秀简书
func UpdateExcellentArticle(id int64, isTop uint8, tx *gorm.DB) error {
	err := tx.Select("is_top").Table(tableNameBook).Where("id = ?", id).Update("is_top", isTop).Error
	return err
}

// UpdateExcellentBlog 管理员设置优秀博客 或者取消
func UpdateExcellentBlog(id int64, isTop uint8, tx *gorm.DB) error {
	err := tx.Select("is_top").Table(tableNameBlog).Where("id = ?", id).Update("is_top", isTop).Error
	return err
}

// SelectArticleByMobile 根据电话查找简书
func SelectArticleByMobile(mobile string) (err error, title string) {
	err = DB.Select("book_title").Table(tableNameBook).Where("mobile = ?", mobile).Take(&title).Error
	return
}

// SelectBlogByMobile 根据电话查找博客
func SelectBlogByMobile(mobile string) (err error, title string) {
	err = DB.Select("book_title").Table(tableNameBlog).Where("mobile = ?", mobile).Take(&title).Error
	return
}

// SelectTopArticle 查询本周的优秀简书
func SelectTopArticle() (err error, excellent []bookBlogArticle.TbBookArticle) {
	err = DB.Where("is_top = 1").Find(&excellent).Error
	return
}

// SelectTopBlog 查询本周博客
func SelectTopBlog() (err error, excellent []bookBlogArticle.TbBlog) {
	err = DB.Where("is_top = 1").Find(&excellent).Error
	return
}

// SelectExcellentCount 查询出优秀简书博客次数排行前5的人
func SelectExcellentCount() (err error, result []usersModel.TbUser) {
	err = DB.Where("excellent_count != 0").
		Where("excellent_count IN (SELECT excellent_count FROM (SELECT excellent_count FROM tb_user WHERE excellent_count != 0 ORDER BY excellent_count DESC LIMIT 5) AS subquery)").
		Find(&result).Error
	//not_written_count
	return
}

// SelectNoWriteCount 查询简书博客未写次数前3的人
func SelectNoWriteCount() (err error, result []usersModel.TbUser) {
	err = DB.Where("not_written_count != 0").
		Where("not_written_count IN (SELECT not_written_count FROM (SELECT not_written_count FROM tb_user WHERE not_written_count != 0 ORDER BY not_written_count DESC LIMIT 5) AS subquery)").
		Find(&result).Error
	return
}
