package main

import (
	"encoding/json"
	"strings"
)

// Получение поля из JSON
// Предусмотрена неограниченная вложенность.
// В случае вложенности поля разделяются точкой, например country.code
func getJsonField(data []byte, field string) string {
	var result string = ""
	var ans map[string]interface{}
	err := json.Unmarshal(data, &ans) // парсим json в карту интерфейсов
	if err != nil {
		return "Wrong JSON"
	}

	// в цикле проходим уровни вложенности и ищем нужное поле
	curMap := &ans
	for len(field) > 0 {
		pointIndex := strings.Index(field, ".")
		var curField string
		if pointIndex >= 0 {
			curField = field[:pointIndex]
		} else {
			curField = field
		}
		value := (*curMap)[curField]
		switch i := value.(type) {
		case string:
			result = i
		case map[string]interface{}:
			curMap = &i
		}
		if pointIndex >= 0 {
			field = field[pointIndex+1:]
		} else {
			field = ""
		}
	}
	return result
}
