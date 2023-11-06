package model

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
