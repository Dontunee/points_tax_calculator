package data

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

type TaxBracketFetcher interface {
	fetchTaxBrackets(year int) (*IncomeTaxCalculatorResponse, error)
}

type HTTPTaxBracketFetcher struct {
	providerURL string
}

type IncomeTaxCalculator struct {
	fetcher TaxBracketFetcher
}

type TaxBracket struct {
	Max  float64 `json:"max,omitempty"` // Omitempty allows Max to be omitted for the highest bracket
	Min  float64 `json:"min"`
	Rate float64 `json:"rate"`
}

type IncomeTaxCalculatorResponse struct {
	TaxBrackets []TaxBracket `json:"tax_brackets"`
}

type IncomeTaxCalculation struct {
	TotalTax      float64                    `json:"totalTax"`
	TaxesPerBand  []IncomeTaxBandCalculation `json:"taxesPerBand"`
	EffectiveRate float64                    `json:"effectiveRate"`
}

type IncomeTaxBandCalculation struct {
	Band      string  `json:"band"`
	TaxedAt   float64 `json:"taxedAt"`
	TaxAmount float64 `json:"taxAmount"`
}

func NewHTTPTaxBracketFetcher(providerURL string) (*HTTPTaxBracketFetcher, error) {
	if providerURL == "" {
		return nil, fmt.Errorf("providerURL is not provided")
	}

	return &HTTPTaxBracketFetcher{providerURL: providerURL}, nil
}

func NewIncomeTaxCalculator(fetcher TaxBracketFetcher) (*IncomeTaxCalculator, error) {
	if fetcher == nil {
		return nil, fmt.Errorf("fetcher is not provided")
	}

	return &IncomeTaxCalculator{fetcher: fetcher}, nil
}

func (tc *IncomeTaxCalculator) CalculateIncomeTax(income float64, year int) (IncomeTaxCalculation, error) {
	var response IncomeTaxCalculation

	taxData, err := tc.fetcher.fetchTaxBrackets(year)
	if err != nil {
		return response, err
	}

	if taxData == nil {
		return response, fmt.Errorf("tax data is nil")
	}

	if len(taxData.TaxBrackets) == 0 {
		return response, fmt.Errorf("tax data has no brackets defined")
	}

	if income < 0 {
		return response, fmt.Errorf("income cannot be negative")
	}

	var totalTax float64
	var taxesPerBand []IncomeTaxBandCalculation

	for i, bracket := range taxData.TaxBrackets {
		if bracket.Min < 0 || bracket.Rate < 0 || (bracket.Max < bracket.Min && bracket.Max != 0) {
			return response, fmt.Errorf("invalid tax bracket data")
		}

		var taxableIncome float64
		var bandMax = bracket.Max

		if i == len(taxData.TaxBrackets)-1 || income < bracket.Max {
			bandMax = income
		}

		if income > bracket.Min {
			taxableIncome = bandMax - bracket.Min
			if taxableIncome < 0 {
				continue
			}

			taxAmount := tc.roundToTwoDecimalPlaces(taxableIncome * bracket.Rate)
			totalTax += tc.roundToTwoDecimalPlaces(taxAmount)

			bandLabel := fmt.Sprintf("%.2f to %.2f", bracket.Min, bracket.Max)
			taxesPerBand = append(taxesPerBand, IncomeTaxBandCalculation{
				Band:      bandLabel,
				TaxedAt:   tc.roundToTwoDecimalPlaces(bracket.Rate * 100), // Converting to percentage
				TaxAmount: taxAmount,
			})
		}

		if income < bracket.Max {
			break
		}
	}

	effectiveRate := 0.0
	if income > 0 {
		effectiveRate = totalTax / income * 100 // Converting to percentage
	}

	// Round totalTax before assigning it to response
	totalTax = tc.roundToTwoDecimalPlaces(totalTax)

	response = IncomeTaxCalculation{
		TotalTax:      totalTax,
		TaxesPerBand:  taxesPerBand,
		EffectiveRate: tc.roundToTwoDecimalPlaces(effectiveRate),
	}

	return response, nil
}

func (tc *IncomeTaxCalculator) roundToTwoDecimalPlaces(number float64) float64 {
	return math.Round(number*100) / 100
}

func (h *HTTPTaxBracketFetcher) fetchTaxBrackets(year int) (*IncomeTaxCalculatorResponse, error) {
	url := fmt.Sprintf("%s/tax-year/%d", h.providerURL, year)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tax bracket fetcher failed with non-OK status code: %d, message: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var taxResponse IncomeTaxCalculatorResponse
	if err := json.Unmarshal(body, &taxResponse); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	if len(taxResponse.TaxBrackets) == 0 {
		return nil, fmt.Errorf("tax data has no brackets defined")
	}

	return &taxResponse, nil
}
