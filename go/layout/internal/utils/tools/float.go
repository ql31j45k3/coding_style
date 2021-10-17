package tools

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

func FloatTruncate(f float64, precision int32) (float64, error) {
	return strconv.ParseFloat(decimal.NewFromFloat(f).Truncate(precision).String(), 64)
}

func FloatTruncateAndPrecisionSecond(f float64) (float64, error) {
	return FloatTruncate(f, 2)
}

func FloatAdd(s float64, p float64) (float64, error) {
	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	return strconv.ParseFloat(f1.Add(f2).Truncate(2).String(), 64)
}

func FloatSub(s float64, p float64) (float64, error) {
	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	return strconv.ParseFloat(f1.Sub(f2).Truncate(2).String(), 64)
}

func FloatDiv(s float64, p float64) (float64, error) {
	return FloatDivAndTruncate(s, p, 4)
}

func FloatDivAndTruncate(s float64, p float64, truncate int32) (float64, error) {
	if p == float64(0) {
		return 0, nil
	}

	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	return strconv.ParseFloat(f1.Div(f2).Truncate(truncate).String(), 64)
}

func FloatMul(s float64, p float64) (float64, error) {
	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	return strconv.ParseFloat(f1.Mul(f2).Truncate(2).String(), 64)
}

func FloatPercent(s float64, p float64) (float64, error) {
	var err error
	var res float64

	res, err = FloatDiv(s, p)
	if err != nil {
		return 0, fmt.Errorf("FloatDiv - %w", err)
	}

	res, err = FloatMul(res, float64(100))
	if err != nil {
		return 0, fmt.Errorf("FloatMul - %w", err)
	}

	if FloatLT(res, 0) {
		res = 0
	}

	return res, nil
}

func FloatGT(s float64, p float64) bool {
	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	return f1.GreaterThan(f2)
}

func FloatLT(s float64, p float64) bool {
	f1 := decimal.NewFromFloat(s)
	f2 := decimal.NewFromFloat(p)

	return f1.LessThan(f2)
}
