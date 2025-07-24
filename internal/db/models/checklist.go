package models

func init() {
	Register("checklists", EntityModel{
		TableName:  "checklists",
		PrimaryKey: "id",
		HasMany: []string{
			"agricultures", "checklist_answers", "checklist_farmer", "checklist_history",
			"checklist_images", "checklist_protocol", "diagnostics", "forestry",
			"livestocks", "scores",
		},
		BelongsTo: map[string]string{
			"analysts:create_by":  "create_by",
			"analysts:aproved_by": "aproved_by",
			"farms":               "farm_id",
			"companies":           "company_id",
		},
		Fields: map[string]FieldDefinition{
			"id":             {SQLType: "int8", GoType: "int64"},
			"create_by":      {SQLType: "int8", GoType: "int64"},
			"aproved_by":     {SQLType: "int8", GoType: "*int64"},
			"farm_id":        {SQLType: "int8", GoType: "*int64"},
			"status_id":      {SQLType: "int8", GoType: "int64"},
			"company_id":     {SQLType: "int8", GoType: "int64"},
			"farm_data":      {SQLType: "json", GoType: "[]byte"},
			"farmers_data":   {SQLType: "json", GoType: "[]byte"},
			"visited_at":     {SQLType: "date", GoType: "time.Time"},
			"was_escorted":   {SQLType: "bool", GoType: "bool"},
			"escort":         {SQLType: "varchar", GoType: "string"},
			"created_at":     {SQLType: "timestamp", GoType: "*time.Time"},
			"updated_at":     {SQLType: "timestamp", GoType: "*time.Time"},
			"processed_pics": {SQLType: "bool", GoType: "bool"},
			"analyst_name":   {SQLType: "varchar", GoType: "*string"},
			"aprt_id":        {SQLType: "int8", GoType: "*int64"},
		},
	})

	Register("analysts", EntityModel{
		TableName:  "analysts",
		PrimaryKey: "id",
		Fields: map[string]FieldDefinition{
			"id":   {SQLType: "int8", GoType: "int64"},
			"name": {SQLType: "varchar", GoType: "string"},
		},
	})

	Register("farms", EntityModel{
		TableName:  "farms",
		PrimaryKey: "id",
		Fields: map[string]FieldDefinition{
			"id":   {SQLType: "int8", GoType: "int64"},
			"name": {SQLType: "varchar", GoType: "string"},
		},
	})

	Register("companies", EntityModel{
		TableName:  "companies",
		PrimaryKey: "id",
		Fields: map[string]FieldDefinition{
			"id":   {SQLType: "int8", GoType: "int64"},
			"name": {SQLType: "varchar", GoType: "string"},
		},
	})

}
