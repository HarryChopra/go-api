package util

// currencies definition for the app
const (
	USD = "USD"
	CAD = "CAD"
	GBP = "GBP"
	EUR = "EUR"
	AUD = "AUD"
)

var supportedCurrencies = []string{USD, CAD, GBP, EUR, AUD}

// IsSupportedCurrency returns true if the currency is supported by the application
func IsSupportedCurrency(currency string) bool {
	for _, supportedCurrency := range supportedCurrencies {
		if supportedCurrency == currency {
			return true
		}
	}
	return false
}
