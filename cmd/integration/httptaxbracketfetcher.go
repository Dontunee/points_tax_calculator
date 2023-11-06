package integration

import (
	"encoding/json"
	"fmt"
	"github.com/Dontunee/points_tax_calculator/cmd/model"
	"io"
	"net/http"
)

type HTTPTaxBracketFetcher struct {
	providerURL string
}

func NewHTTPTaxBracketFetcher(providerURL string) (*HTTPTaxBracketFetcher, error) {
	if providerURL == "" {
		return nil, fmt.Errorf("providerURL is not provided")
	}

	return &HTTPTaxBracketFetcher{providerURL: providerURL}, nil
}

func (h *HTTPTaxBracketFetcher) FetchTaxBrackets(year int) (*model.IncomeTaxCalculatorResponse, error) {
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

	var taxResponse model.IncomeTaxCalculatorResponse
	if err := json.Unmarshal(body, &taxResponse); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	if len(taxResponse.TaxBrackets) == 0 {
		return nil, fmt.Errorf("tax data has no brackets defined")
	}

	return &taxResponse, nil
}
