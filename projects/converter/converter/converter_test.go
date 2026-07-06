package converter

import (
	"math"
	"testing"
)

func TestParseCurrency(t *testing.T) {
	tests := []CurrencyCode{
		RUB,
		USD,
		EUR,
		CNY,
	}

	for _, tt := range tests {
		t.Run(tt.String(), func(t *testing.T) {
			result := ParseCurrency(tt.String())
			if result != tt {
				t.Errorf(`Invalid currency, expected: %s, actual: %s`, tt.String(), result)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	result, err := Convert(100, USD, RUB)

	if err != nil {
		t.Errorf("error in convertion, %s", err)
	}

	expected := 8140.0 // 100 * 81.4
	if !almostEqual(result, expected) {
		t.Errorf("expected: %f, actual: %f", result, expected)
	}
}

func TestConvertNegativeAmount(t *testing.T) {
	_, err := Convert(-1, USD, RUB)

	if err == nil {
		t.Errorf("Expected error, amount = -1")
	}
}

const epsilon = 1e-9 // допустимая погрешность

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}
