package main

import (
	"flag"
	"fmt"
	"projects/kim/converter/converter"
)

func main() {
	amount, from, to := initFlags()

	toAmount, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Println("Error in convertion", err)
		return
	}
	fmt.Printf("Converted %d %s to %.2f %s", amount, from.String(), toAmount, to.String())
}

func initFlags() (int, converter.CurrencyCode, converter.CurrencyCode) {

	var fromCurrency = flag.String("from-currency", "usd", "from currency")
	var toCurrency = flag.String("to-currency", "rub", "to currency")
	var amount = flag.Int("amount", 1000, "amount")
	flag.Parse()

	from := converter.ParseCurrency(*fromCurrency)
	to := converter.ParseCurrency(*toCurrency)

	return *amount, from, to
}
