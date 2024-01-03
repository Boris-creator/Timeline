package validators

import (
	"fiber-server/db"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Exists(fl validator.FieldLevel) bool {
	//if !fl.Field().IsValid() {
	//return true
	//}
	//return true

	params := strings.Split(fl.Param(), ".")
	table := params[0]
	column := "id"
	if len(params) > 1 {
		column = params[1]
	}

	var value any
	switch valueType := fl.Field().Kind(); valueType {
	case reflect.String:
		value = fl.Field().String()
	case reflect.Int:
		value = fl.Field().Int()
	}

	var result int
	db.Database.Raw(
		fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column),
		value,
	).Scan(&result)

	return result > 0
}
