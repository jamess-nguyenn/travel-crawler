package helper

func CalculateOffset(page int, pageSize int) int {
	return (page - 1) * pageSize
}
