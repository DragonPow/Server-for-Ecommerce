package math

import (
	"fmt"
	"golang.org/x/exp/maps"
	"math"
)

const (
	ZERO_INT                  = 0
	DEFAULT_PRECISION float64 = 0.0001
	BaseTwo           float64 = 2.0
	Phase53           float64 = 53.0
)

// RoundFloat
//  val: the value want round
//  roundingFactor: must be 10^-x, ex: 0.1, 0.001, 0.0001,...
//  output: the val after ROUND_UP
//  Example:
//   (1.4447, 0.001) => 1.445
//   (3.00000123, 0.1) => 3.0
//   (4.9999, 0.0001) => 5.0
//   (2.34, 0.01) => 2.34
//   (13.5, 1) => 14
//   (167.5, 10) => 168 // Warning: Should not use this case, this func must only use for round decimal in float
func RoundFloat(val, roundingFactor float64) float64 {
	if roundingFactor == ZERO_INT || val == ZERO_INT {
		return ZERO_INT
	}
	precision := math.Ceil(math.Abs(math.Log10(roundingFactor)))
	pow10Precision := math.Pow10(int(precision))
	return math.Round(val*pow10Precision) / pow10Precision
}

func ToMap[T, U any, K comparable](slice []T, f func(item T) (K, U)) map[K]U {
	result := make(map[K]U)
	for _, i := range slice {
		k, v := f(i)
		result[k] = v
	}
	return result
}

func Convert[T, U any](slice []T, f func(item T) U) []U {
	result := make([]U, 0, len(slice))
	for _, i := range slice {
		v := f(i)
		result = append(result, v)
	}
	return result
}

func ConvertMap[T, Z comparable, U, Y any](myMap map[T]U, f func(key T, value U) (Z, Y)) map[Z]Y {
	result := make(map[Z]Y, len(myMap))
	for k, v := range myMap {
		newK, newV := f(k, v)
		result[newK] = newV
	}
	return result
}

func AppendMap[T comparable, U any](anotherMaps ...map[T]U) (newMap map[T]U) {
	sumLen := 0
	for _, anotherMap := range anotherMaps {
		sumLen += len(anotherMap)
	}
	newMap = make(map[T]U, sumLen)
	for _, anotherMap := range anotherMaps {
		for k, v := range anotherMap {
			_, ok := newMap[k]
			if ok {
				panic(fmt.Sprintf("Key %v is duplicate", k))
			}
			newMap[k] = v
		}
	}
	return newMap
}

func Uniq[T comparable](slice []T) []T {
	result := make(map[T]struct{}, len(slice))
	for _, v := range slice {
		result[v] = struct{}{}
	}
	return maps.Keys(result)
}
