package tools

import (
	"errors"
)

func GetPageStartAndEnd(total, page, limit int) (int, int, error) {
	skipCount := GetSkipCount(page, limit)

	lastResultCount := total - skipCount
	if lastResultCount < 0 {
		return 0, 0, errors.New("lastResultCount < 0, Should not be less than 0")
	}

	sliceStart := skipCount
	sliceEnd := 0

	if lastResultCount < limit {
		sliceEnd = skipCount + lastResultCount
	} else {
		sliceEnd = skipCount + limit
	}

	return sliceStart, sliceEnd, nil
}

func GetSkipCount(page, limit int) int {
	return (page - 1) * limit
}
