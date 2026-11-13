package airport

const defaultSortColumn = tableName + ".id"

var sortColumns = map[string]string{
	"id":         defaultSortColumn,
	"iata_code":  tableName + ".iata_code," + defaultSortColumn,
	"city":       tableName + ".city," + defaultSortColumn,
	"country":    tableName + ".country," + defaultSortColumn,
	"created_at": tableName + ".created_at," + defaultSortColumn,
	"created_by": tableName + ".created_by," + defaultSortColumn,
	"updated_at": tableName + ".updated_at," + defaultSortColumn,
	"updated_by": tableName + ".updated_by," + defaultSortColumn,
}
