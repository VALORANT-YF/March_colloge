package mysql

import (
	"college/models/bookBlogArticle"
	"college/models/deptsModel"
	"college/models/robotModels"
	"college/models/usersModel"
	"college/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 初始化MySQL连接r
func Init(cfg *settings.MySQLConfig) (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	DB, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{})
	if err != nil {
		return
	} else {
		sqlDB, _ := DB.DB()
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConnection)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConnection)
	}

	//自动建表
	DB.AutoMigrate(bookBlogArticle.TbBookArticle{})
	DB.AutoMigrate(bookBlogArticle.TbBlog{})
	DB.AutoMigrate(deptsModel.TbDept{})
	DB.AutoMigrate(usersModel.TbUser{})
	DB.AutoMigrate(robotModels.TbRobot{})
	//migrate 仅支持创建表、增加表中没有的字段和索引
	//DB.AutoMigrate(&models.Student{})
	return
}

// 自动建表方法
func creatTable(dst interface{}) {
	if !DB.Migrator().HasTable(dst) {
		err := DB.AutoMigrate(dst)
		if err != nil {
			return
		}
		if DB.Migrator().HasTable(dst) {
			fmt.Println("表创建成功")
		} else {
			fmt.Println("表创建失败")
		}
	} else {
		fmt.Println("表已存在")
	}
}
