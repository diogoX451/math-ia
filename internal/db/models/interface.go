package models

type EntityModel struct {
	TableName  string
	PrimaryKey string
	HasMany    []string
	BelongsTo  map[string]string
	Fields     map[string]FieldDefinition
}

type FieldDefinition struct {
	SQLType string
	GoType  string
}

var Registry = map[string]EntityModel{}
var TableToEntity = map[string]string{}
var Fields = map[string]map[string]FieldDefinition{}

func Register(name string, model EntityModel) {
	Registry[name] = model
	TableToEntity[model.TableName] = name
	Fields[name] = model.Fields
}
