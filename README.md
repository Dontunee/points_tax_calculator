# Points Tax Calculator

This repository contains an HTTP API that calculates total taxes owed based on annual income and tax year. The API endpoint returns detailed tax information, including the amount of taxes owed per band and the effective tax rate.

## Prerequisites

Before starting, ensure you have the following installed:
* [Go programming language](https://golang.org/dl/)

## Installation

To install, follow these steps:
1. Clone this repository.
2. Run `make download` in the terminal to install dependencies.

## Running Locally

To run the API locally:
- follow the steps [here](https://github.com/points/interview-test-server#get-up-and-running) to run the tax bracket provider server(This is required for this to run locally)
- Use the command `make run/api` in the terminal.
- The API will start on port 8000 by default.

## Running tests
  To run unit tests:
- Use the command `make test` in the terminal.


## Endpoint Documentation

To use the API:
- Make a request to `{serverUrl}/v1/tax-calculator?income={incomeValue}&taxYear={taxYear}`.
- You will receive a JSON response like this: `{"totalTax":17739.17,"taxesPerBand":[{"band":"0.00 to 50197.00","taxedAt":15,"taxAmount":7529.55},{"band":"50197.00 to 100392.00","taxedAt":20.5,"taxAmount":10209.62}],"effectiveRate":17.74}`.

The endpoint currently supports the tax years 2019, 2020, 2021, and 2022.

## Known Issues
- Swagger/OpenAPI documentation not available yet.
- Rate limiting not implemented.
- API authentication not set up.
