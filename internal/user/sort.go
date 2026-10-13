package user

const defaultSortColumn = tableName + ".id"

var sortColumns = map[string]string{
	"id":                defaultSortColumn,
	"email":             tableName + ".email," + defaultSortColumn,
	"auth0_id":          tableName + ".auth0_id," + defaultSortColumn,
	"email_verified_at": tableName + ".email_verified_at," + defaultSortColumn,
	"first_name":        tableName + ".first_name," + defaultSortColumn,
	"last_name":         tableName + ".last_name," + defaultSortColumn,
	"full_name":         tableName + ".full_name," + defaultSortColumn,
	"phone_number":      tableName + ".phone_number," + defaultSortColumn,
	"status":            tableName + ".status," + defaultSortColumn,
	"last_login_at":     tableName + ".last_login_at," + defaultSortColumn,
	"created_at":        tableName + ".created_at," + defaultSortColumn,
	"created_by":        tableName + ".created_by," + defaultSortColumn,
	"updated_at":        tableName + ".updated_at," + defaultSortColumn,
	"updated_by":        tableName + ".updated_by," + defaultSortColumn,
}
