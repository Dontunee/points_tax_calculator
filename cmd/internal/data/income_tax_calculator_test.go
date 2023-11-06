package data

import (
	"fmt"
	"github.com/Dontunee/points_tax_calculator/cmd/model"
	"testing"
)

type MockTaxBracketFetcher struct {
	MockResponse *model.IncomeTaxCalculatorResponse
	MockError    error
}

func (m *MockTaxBracketFetcher) FetchTaxBrackets(year int) (*model.IncomeTaxCalculatorResponse, error) {
	return m.MockResponse, m.MockError
}

func TestRoundToTwoDecimalPlaces(t *testing.T) {
	tests := []struct {
		name   string
		number float64
		want   float64
	}{
		{"Round Down", 0.123, 0.12},
		{"Round Up", 0.125, 0.13},
		{"No Rounding", 0.1, 0.1},
	}

	calculator := &IncomeTaxCalculator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculator.roundToTwoDecimalPlaces(tt.number); got != tt.want {
				t.Errorf("roundToTwoDecimalPlaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateIncomeTax(t *testing.T) {
	// Mock successful response for tax brackets
	successfulMockResponse := &model.IncomeTaxCalculatorResponse{
		TaxBrackets: []model.TaxBracket{
			{Min: 0, Max: 50197, Rate: 0.15},
			{Min: 50197, Max: 100392, Rate: 0.205},
			{Min: 100392, Max: 155625, Rate: 0.26},
			{Min: 155625, Max: 221708, Rate: 0.29},
			{Min: 221708, Rate: 0.33},
		},
	}

	tests := []struct {
		name              string
		mockFetcher       *MockTaxBracketFetcher
		income            float64
		year              int
		wantTotalTax      float64
		wantEffectiveRate float64
		wantErr           bool
	}{
		{"Valid Income Low Bracket", &MockTaxBracketFetcher{MockResponse: successfulMockResponse}, 0, 2022, 0, 0, false},
		{"Valid Income Middle Bracket", &MockTaxBracketFetcher{MockResponse: successfulMockResponse}, 50000, 2022, 7500, 15, false},
		{"Valid Income High Bracket", &MockTaxBracketFetcher{MockResponse: successfulMockResponse}, 100000, 2022, 17739.17, 17.74, false},
		{"Negative Income", &MockTaxBracketFetcher{MockResponse: successfulMockResponse}, 1234567, 2022, 385587.64, 31.23, false},
		{"Negative Income", &MockTaxBracketFetcher{MockResponse: successfulMockResponse}, -1000, 2022, 0, 0, true},
		{"Fetch Tax Bracket Failure", &MockTaxBracketFetcher{MockError: fmt.Errorf("fetch error")}, 50000, 2022, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := &IncomeTaxCalculator{fetcher: tt.mockFetcher}
			got, err := calculator.CalculateIncomeTax(tt.income, tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateIncomeTax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.TotalTax != tt.wantTotalTax {
					t.Errorf("CalculateIncomeTax() got total tax = %v, want %v", got.TotalTax, tt.wantTotalTax)
				}
				if got.EffectiveRate != tt.wantEffectiveRate {
					t.Errorf("CalculateIncomeTax() got effective rate = %v, want %v", got.EffectiveRate, tt.wantEffectiveRate)
				}
			}
		})
	}
}
