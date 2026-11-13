package airline

const defaultSortColumn = tableName + ".id"

var sortColumns = map[string]string{
	"id":            defaultSortColumn,
	"provider_code": tableName + ".provider_code," + defaultSortColumn,
	"iata_code":     tableName + ".iata_code," + defaultSortColumn,
	"name":          tableName + ".name," + defaultSortColumn,
	"alliance":      tableName + ".alliance," + defaultSortColumn,
	"active":        tableName + ".active," + defaultSortColumn,
	"created_at":    tableName + ".created_at," + defaultSortColumn,
	"created_by":    tableName + ".created_by," + defaultSortColumn,
	"updated_at":    tableName + ".updated_at," + defaultSortColumn,
	"updated_by":    tableName + ".updated_by," + defaultSortColumn,
}
