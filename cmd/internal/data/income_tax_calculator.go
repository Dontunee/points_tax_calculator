package data

import (
	"fmt"
	"github.com/Dontunee/points_tax_calculator/cmd/model"
	"math"
)

type IncomeTaxCalculator struct {
	fetcher model.TaxBracketFetcher
}

func NewIncomeTaxCalculator(fetcher model.TaxBracketFetcher) (*IncomeTaxCalculator, error) {
	if fetcher == nil {
		return nil, fmt.Errorf("fetcher is not provided")
	}

	return &IncomeTaxCalculator{fetcher: fetcher}, nil
}

func (tc *IncomeTaxCalculator) CalculateIncomeTax(income float64, year int) (model.IncomeTaxCalculation, error) {
	var response model.IncomeTaxCalculation

	taxData, err := tc.fetcher.FetchTaxBrackets(year)
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
	var taxesPerBand []model.IncomeTaxBandCalculation

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
			taxesPerBand = append(taxesPerBand, model.IncomeTaxBandCalculation{
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

	totalTax = tc.roundToTwoDecimalPlaces(totalTax)

	response = model.IncomeTaxCalculation{
		TotalTax:      totalTax,
		TaxesPerBand:  taxesPerBand,
		EffectiveRate: tc.roundToTwoDecimalPlaces(effectiveRate),
	}

	return response, nil
}

func (tc *IncomeTaxCalculator) roundToTwoDecimalPlaces(number float64) float64 {
	return math.Round(number*100) / 100
}
