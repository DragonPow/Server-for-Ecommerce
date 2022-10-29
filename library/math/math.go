package math

import "math"

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
