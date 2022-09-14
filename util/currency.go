package util

var supportedCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"INR": true,
	"CAD": true,
	"YEN": true,
}

func IsSupportedCurrency(currency string) bool {
	return supportedCurrencies[currency]
}
