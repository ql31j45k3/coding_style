package tools

import (
	"strconv"
	"strings"
)

func IsEmpty(str string) bool {
	str = strings.TrimSpace(str)
	return str == ""
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

func Atoi(str string, defaultValue int) (int, error) {
	if IsEmpty(str) {
		return defaultValue, nil
	}

	result, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue, err
	}
	return result, nil
}

// ArrayStrUnique 判斷陣列中的值是否重複
func ArrayStrUnique(list []string) []string {
	keys := map[string]bool{}
	strList := make([]string, 0)

	for _, element := range list {
		if _, value := keys[element]; !value {
			keys[element] = true
			strList = append(strList, element)
		}
	}
	return strList
}

func ArrayFind(list []string, subStr string) bool {
	for i := range list {
		if list[i] == subStr {
			return true
		}
	}

	return false
}
