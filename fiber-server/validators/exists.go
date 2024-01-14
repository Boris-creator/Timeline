package validators

import (
	"fiber-server/db"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Exists(fl validator.FieldLevel) bool {
	table, column, value := getRequestParams(fl)

	return existsInTable(table, column, value)
}

func Unique(fl validator.FieldLevel) bool {
	table, column, value := getRequestParams(fl)

	return !existsInTable(table, column, value)
}

func existsInTable(tableName string, columnName string, value any) bool {
	var result int
	db.Database.Raw(
		fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", tableName, columnName),
		value,
	).Scan(&result)

	return result > 0
}

func getRequestParams(fl validator.FieldLevel) (table, column string, value any) {
	params := strings.Split(fl.Param(), ".")
	table = params[0]
	column = "id"
	if len(params) > 1 {
		column = params[1]
	}

	switch valueType := fl.Field().Kind(); valueType {
	case reflect.String:
		value = fl.Field().String()
	case reflect.Int:
		value = fl.Field().Int()
	}

	return table, column, value
}
