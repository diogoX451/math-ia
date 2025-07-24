package loader

import (
	"fmt"
	"math-ia/internal/db"
	"math-ia/internal/db/models"
	"math-ia/internal/tools/formater"
	"strings"
)

func GetEntityContext(entityKey string, id int64) ([]string, error) {
	var context []string

	model, ok := models.Registry[entityKey]
	if !ok {
		return context, nil
	}

	row := db.GetDB().QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", model.TableName, model.PrimaryKey), id)
	if text, _ := formater.FormatRow(model.TableName, row); text != "" {
		context = append(context, fmt.Sprintf("%s: %s", model.TableName, text))
	}

	for _, relatedTable := range model.HasMany {
		q := fmt.Sprintf("SELECT * FROM %s WHERE %s_id = $1 LIMIT 3", relatedTable, model.TableName)
		rows, err := db.GetDB().Query(q, id)
		if err == nil {
			defer rows.Close()
			context = append(context, formater.FormatRows(relatedTable, rows)...)
		}
	}

	for relatedEntityKey, localField := range model.BelongsTo {
		relatedModel, ok := models.Registry[relatedEntityKey]
		if !ok {
			continue
		}

		q := fmt.Sprintf(
			"SELECT * FROM %s WHERE %s = (SELECT %s FROM %s WHERE %s = $1)",
			relatedModel.TableName,
			relatedModel.PrimaryKey,
			localField,
			model.TableName,
			model.PrimaryKey,
		)
		row := db.GetDB().QueryRow(q, id)
		if text, _ := formater.FormatRow(relatedModel.TableName, row); text != "" {
			context = append(context, fmt.Sprintf("%s: %s", relatedModel.TableName, text))
		}
	}

	return context, nil
}

func GetAllIDs(entityKey string) ([]int64, error) {
	model, ok := models.Registry[entityKey]
	if !ok {
		return nil, fmt.Errorf("entity %s not found", entityKey)
	}

	q := fmt.Sprintf("SELECT %s FROM %s", model.PrimaryKey, model.TableName)
	rows, err := db.GetDB().Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func GetLast10IDs(entity string) ([]int64, error) {
	model, ok := models.Registry[entity]
	if !ok {
		return nil, fmt.Errorf("entidade %s não registrada", entity)
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s ORDER BY %s DESC LIMIT 10",
		model.PrimaryKey, model.TableName, model.PrimaryKey,
	)

	rows, err := db.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func GetRowAsText(entity string, id int64) (string, error) {
	model, ok := models.Registry[entity]
	if !ok {
		return "", fmt.Errorf("entidade %s não registrada", entity)
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", model.TableName, model.PrimaryKey)
	row := db.GetDB().QueryRow(query, id)
	text, err := formater.FormatRow(model.TableName, row)
	return text, err
}

func ContextRecursivo(entity string, id int64, visited map[string]bool) ([]string, error) {
	var result []string
	key := fmt.Sprintf("%s:%d", entity, id)
	if visited[key] {
		return result, nil
	}

	visited[key] = true

	model, ok := models.Registry[entity]

	if !ok {
		return nil, fmt.Errorf("entidade %s não registrada", entity)
	}

	mainRow, err := GetRowAsText(entity, id)
	if err == nil && mainRow != "" {
		result = append(result, fmt.Sprintf("%s [%d]: %s", entity, id, mainRow))
	}

	for targetKey, fk := range model.BelongsTo {
		parts := strings.Split(targetKey, ":")
		entityName := parts[0]

		var fkID int64
		err := db.GetDB().QueryRow(fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", fk, model.TableName, model.PrimaryKey), id).Scan(&fkID)
		if err == nil {
			realEntity, found := models.TableToEntity[entityName]
			if !found {
				fmt.Printf("Entidade %s não encontrada no TableToEntity\n", entityName)
				continue
			}
			relCtx, err := ContextRecursivo(realEntity, fkID, visited)
			if err == nil {
				result = append(result, relCtx...)
			}
		}
	}

	for _, childTable := range model.HasMany {
		rows, err := db.GetDB().Query(fmt.Sprintf("SELECT id FROM %s WHERE %s = $1", childTable, entity+"_id"), id)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var childID int64
				if err := rows.Scan(&childID); err == nil {
					targetEntity, found := models.TableToEntity[childTable]
					if !found {
						fmt.Printf("Tabela %s não registrada no TableToEntity\n", childTable)
						continue
					}
					relCtx, err := ContextRecursivo(targetEntity, childID, visited)
					if err == nil {
						result = append(result, relCtx...)
					}
				}
			}
		}
	}

	return result, nil
}
