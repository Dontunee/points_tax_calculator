package integration

import (
	"encoding/json"
	"github.com/Dontunee/points_tax_calculator/cmd/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHTTPTaxBracketFetcher(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid URL", "http://example.com", false},
		{"Empty URL", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewHTTPTaxBracketFetcher(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIncomeTaxCalculator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchTaxBrackets(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*httptest.Server)
		wantErr   bool
	}{
		{
			name: "TestSuccessfulFetch",
			setupMock: func(server *httptest.Server) {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(model.IncomeTaxCalculatorResponse{
						TaxBrackets: []model.TaxBracket{
							{Min: 0, Max: 10000, Rate: 0.1},
							{Min: 10001, Max: 20000, Rate: 0.2},
						},
					})
				})
			},
			wantErr: false,
		},
		{
			name: "TestErrorOnHTTPRequest",
			setupMock: func(server *httptest.Server) {
				server.Close()
			},
			wantErr: true,
		},
		{
			name: "TestNonOKHTTPStatus",
			setupMock: func(server *httptest.Server) {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				})
			},
			wantErr: true,
		},
		{
			name: "TestErrorReadingResponseBody",
			setupMock: func(server *httptest.Server) {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte{})
				})
			},
			wantErr: true,
		},
		{
			name: "TestInvalidJSONResponse",
			setupMock: func(server *httptest.Server) {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("invalid JSON"))
				})
			},
			wantErr: true,
		},
		{
			name: "TestEmptyJSONResponse",
			setupMock: func(server *httptest.Server) {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{}"))
				})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			defer mockServer.Close()

			// Apply the specific setup for the test case
			tt.setupMock(mockServer)

			taxFetcher, _ := NewHTTPTaxBracketFetcher(mockServer.URL)

			_, err := taxFetcher.FetchTaxBrackets(2022)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
