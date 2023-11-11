package mysql

// UpdateAdmin 设置管理员
func UpdateAdmin(isBoss int, userid string) error {
	return DB.Where("userid = ?", userid).Table("tb_user").Update("is_boss", isBoss).Error
}
