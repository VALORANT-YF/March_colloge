package deptsModel

type DeptIsWrite struct {
	DeptId       int64 `json:"dept_id,omitempty" form:"dept_id"`
	IsWriteBooks uint8 `json:"is_write_books,omitempty" form:"is_write_books"`
}
