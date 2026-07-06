package converter

import (
	"errors"
	"strings"
)

var InvalidAmount = errors.New("amount cannot be zero or negative")

type CurrencyCode int

const (
	RUB CurrencyCode = iota // 0
	USD                     // 1
	EUR                     // 2
	CNY                     // 3
)

func (c CurrencyCode) String() string {
	names := [...]string{"RUB", "USD", "EUR", "CNY"}
	if int(c) < 0 || int(c) >= len(names) {
		return "Unknown"
	}
	return names[c]
}

type Currency struct {
	Rate float64
	Code CurrencyCode
}

func ParseCurrency(currency string) CurrencyCode {
	mapping := map[string]CurrencyCode{
		"rub": RUB,
		"usd": USD,
		"eur": EUR,
		"cny": CNY,
	}
	return mapping[strings.ToLower(currency)]
}

func Convert(amount int, from, to CurrencyCode) (float64, error) {
	if amount <= 0 {
		return -1, InvalidAmount
	}

	rates := RateMap()
	fromRate := rates[from]
	toRate := rates[to]
	rubAmount := float64(amount) * fromRate
	return (rubAmount / toRate), nil
}

func RateMap() map[CurrencyCode]float64 {
	return map[CurrencyCode]float64{
		RUB: 1,
		USD: 81.4,
		EUR: 89.8,
		CNY: 14.91,
	}
}
