package accommodation

const defaultSortColumn = tableName + ".id"

var sortColumns = map[string]string{
	"id":                 defaultSortColumn,
	"provider_code":      tableName + ".provider_code," + defaultSortColumn,
	"external_id":        tableName + ".external_id," + defaultSortColumn,
	"external_parent_id": tableName + ".external_parent_id," + defaultSortColumn,
	"name":               tableName + ".name," + defaultSortColumn,
	"property_type":      tableName + ".property_type," + defaultSortColumn,
	"chain":              tableName + ".chain," + defaultSortColumn,
	"currency":           tableName + ".currency," + defaultSortColumn,
	"country":            tableName + ".country," + defaultSortColumn,
	"city":               tableName + ".city," + defaultSortColumn,
	"latitude":           tableName + ".latitude," + defaultSortColumn,
	"longitude":          tableName + ".longitude," + defaultSortColumn,
	"address":            tableName + ".address," + defaultSortColumn,
	"zip":                tableName + ".zip," + defaultSortColumn,
	"stars":              tableName + ".stars," + defaultSortColumn,
	"rating":             tableName + ".rating," + defaultSortColumn,
	"review_count":       tableName + ".review_count," + defaultSortColumn,
	"created_at":         tableName + ".created_at," + defaultSortColumn,
	"created_by":         tableName + ".created_by," + defaultSortColumn,
	"updated_at":         tableName + ".updated_at," + defaultSortColumn,
	"updated_by":         tableName + ".updated_by," + defaultSortColumn,
}
