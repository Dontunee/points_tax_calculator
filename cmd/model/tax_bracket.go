package model

type TaxBracket struct {
	Max  float64 `json:"max,omitempty"` // Omitempty allows Max to be omitted for the highest bracket
	Min  float64 `json:"min"`
	Rate float64 `json:"rate"`
}

type TaxBracketFetcher interface {
	FetchTaxBrackets(year int) (*IncomeTaxCalculatorResponse, error)
}

type IncomeTaxCalculatorResponse struct {
	TaxBrackets []TaxBracket `json:"tax_brackets"`
}
