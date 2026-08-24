package filter

import (
	"regexp"

	"github.com/scalent.io/scalent-hrms/entity/filters"
	"github.com/scalent.io/scalent-hrms/model"
)

const NO_OF_VALUES_FOR_BETWEEN_CONDITION = 2

var FilterConditionsMap = map[string][]string{
	"int":       {"btw", "gt", "lt", "eq", "in", "like", "notin"},
	"bigint":    {"btw", "gt", "lt", "eq"},
	"string":    {"eq", "like"},
	"varchar":   {"eq", "like", "notin", "in"},
	"char":      {"eq", "like"},
	"text":      {"eq", "like"},
	"datetime":  {"btw", "gt", "lt", "eq", "gteq", "lteq"},
	"date":      {"btw", "gt", "lt", "eq", "gteq", "lteq"},
	"timestamp": {"btw", "gt", "lt", "eq"},
	"float":     {"btw", "gt", "lt", "eq"},
	"double":    {"btw", "gt", "lt", "eq"},
	"decimal":   {"btw", "gt", "lt", "eq"},
	"enum":      {"eq", "in", "notin"},
	"json":      {"JSON_SEARCH"},
}

var SortTypeString = []string{"ASC", "asc", "DESC", "desc"}

func CreateFilterStr(filterStruct []filters.Filter, modelMAP map[string]model.FieldStruct) ([]string, []interface{}) {
	argmnt := []interface{}{}
	whereStr := []string{}

	var fieldStr model.FieldStruct
	var conditions []string
	var ok bool

	for _, filter := range filterStruct {

		if fieldStr, ok = modelMAP[filter.Field]; !ok {
			continue
		}

		if conditions, ok = FilterConditionsMap[fieldStr.MySQLDatatype]; !ok {
			continue
		}

		if conditionMatched := Contains(conditions, filter.Condition); !conditionMatched {
			continue
		}

		str, args := GetWhereStr(filter.Condition, len(filter.FilterValues), fieldStr.FieldName, fieldStr.MySQLDatatype, filter.FilterValues)
		whereStr = append(whereStr, str)
		argmnt = append(argmnt, args...)

	}

	return whereStr, argmnt
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetWhereStr(filterCondition string, valArrayLength int, fieldName string, fieldDataType string, filterValues []string) (string, []interface{}) {

	var str string
	args := []interface{}{}

	if valArrayLength == 0 {
		return "", nil
	}

	if filterCondition == "btw" {
		if NO_OF_VALUES_FOR_BETWEEN_CONDITION == valArrayLength {
			if fieldDataType == "datetime" || fieldDataType == "timestamp" {
				if matched := CheckTimeAvailableInDate(filterValues[0]); !matched {
					filterValues[0] = filterValues[0] + " 00:00:00"
				}
				if matched := CheckTimeAvailableInDate(filterValues[1]); !matched {
					filterValues[1] = filterValues[1] + " 23:59:59"
				}
			}
			args = append(args, filterValues[0], filterValues[1])

			str = fieldName + " BETWEEN ? AND ? "

		}

	} else if filterCondition == "eq" {
		str = fieldName + " = ?"
		args = append(args, filterValues[0])

	} else if filterCondition == "gt" {
		args = append(args, filterValues[0])
		str = fieldName + " > ?"

	} else if filterCondition == "lt" {
		args = append(args, filterValues[0])
		str = fieldName + " < ?"

	} else if filterCondition == "gteq" {
		args = append(args, filterValues[0])
		str = fieldName + " >= ?"

	} else if filterCondition == "lteq" {
		args = append(args, filterValues[0])
		str = fieldName + " <= ?"

	} else if filterCondition == "like" {
		str = fieldName + " LIKE ?"
		likeValStr := `%` + filterValues[0] + `%`
		args = append(args, likeValStr)

	} else if filterCondition == "in" {
		var qMark string
		for i := 0; i < valArrayLength; i++ {
			if i == valArrayLength-1 {
				qMark = qMark + "?"
				args = append(args, filterValues[i])
			} else {
				qMark = qMark + "?, "
				args = append(args, filterValues[i])
			}
		}
		str = fieldName + " IN(" + qMark + ")"
	} else if filterCondition == "notin" {
		var qMark string
		for i := 0; i < valArrayLength; i++ {
			if i == valArrayLength-1 {
				qMark = qMark + "?"
				args = append(args, filterValues[i])
			} else {
				qMark = qMark + "?, "
				args = append(args, filterValues[i])
			}
		}
		str = fieldName + " NOT IN(" + qMark + ")"
	} else if filterCondition == "JSON_SEARCH" {
		str = `JSON_SEARCH(tags, 'all', ?, NULL, '$[*]') IS NOT NULL`
		searchStr := `%` + filterValues[0] + `%`
		args = append(args, searchStr)
	}

	return str, args
}

func CheckTimeAvailableInDate(date string) bool {
	re := regexp.MustCompile(`\d{2}:\d{2}:\d{2}`)

	regEx := re.Match([]byte(date))
	if !regEx {
		return false
	} else {
		return true
	}

}

func CreateSortStr(sortStruct filters.SortOption, modelMAP map[string]model.FieldStruct) string {
	var str string
	var val model.FieldStruct
	var ok, conditionMatched bool

	val, ok = modelMAP[sortStruct.SortBy]
	if ok {
		conditionMatched = Contains(SortTypeString, sortStruct.SortType)
		if conditionMatched {
			str = " ORDER BY " + val.FieldName + " " + sortStruct.SortType
		}
	}

	return str
}
