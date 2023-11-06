package main

import (
	"fmt"
	"github.com/Dontunee/points_tax_calculator/cmd/internal/data"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type IncomeTaxCalculationResponse struct {
	TotalTax      float64                      `json:"totalTax"`
	TaxesPerBand  []TaxBandCalculationResponse `json:"taxesPerBand"`
	EffectiveRate float64                      `json:"effectiveRate"`
}

type TaxBandCalculationResponse struct {
	Band      string  `json:"band"`
	TaxedAt   float64 `json:"taxedAt"`
	TaxAmount float64 `json:"taxAmount"`
}

func (app *application) calculateIncomeTaxHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parsing query parameters
	query := r.URL.Query()
	incomeStr := query.Get("income")
	taxYearStr := query.Get("taxYear")

	if incomeStr == "" {
		app.errorResponse(w, http.StatusBadRequest, fmt.Errorf("missing required parameter: income"))
		return
	}

	if taxYearStr == "" {
		app.errorResponse(w, http.StatusBadRequest, fmt.Errorf("missing required parameter: taxYear"))
		return
	}

	income, err := strconv.ParseFloat(incomeStr, 64)
	if err != nil {
		app.errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid income parameter"))
		return
	}

	taxYear, err := strconv.Atoi(taxYearStr)
	if err != nil {
		app.errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid tax year parameter"))
		return
	}

	httpTaxBracketFetcher, err := data.NewHTTPTaxBracketFetcher(app.config.taxCalculatorUrl)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	taxCalculator, err := data.NewIncomeTaxCalculator(httpTaxBracketFetcher)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	tax, err := taxCalculator.CalculateIncomeTax(income, taxYear)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	taxCalculationResponse := copyTaxCalculation(tax)
	app.writeJSON(w, 200, taxCalculationResponse, nil)
	return
}

func copyTaxCalculation(src data.IncomeTaxCalculation) IncomeTaxCalculationResponse {
	var dst IncomeTaxCalculationResponse

	dst.TotalTax = src.TotalTax
	dst.EffectiveRate = src.EffectiveRate

	// Deep copy of slice of IncomeTaxBandCalculation
	for _, band := range src.TaxesPerBand {
		copiedBand := TaxBandCalculationResponse{
			Band:      band.Band,
			TaxedAt:   band.TaxedAt,
			TaxAmount: band.TaxAmount,
		}
		dst.TaxesPerBand = append(dst.TaxesPerBand, copiedBand)
	}

	return dst
}
