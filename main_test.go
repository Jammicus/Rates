package main

import (
	"testing"
)

func TestGenerateRequest(t *testing.T) {
	var testcases = []struct {
		command   string
		baseFlag  string
		startFlag string
		endFlag   string
		currency  string
		expected  string
	}{
		{"latest", "", "", "", "", "https://api.exchangeratesapi.io/latest"},
		{"latest", "GBP", "", "", "", "https://api.exchangeratesapi.io/latest?base=GBP"},
		{"latest", "", "", "", "gbp,usd", "https://api.exchangeratesapi.io/latest?symbols=GBP,USD"},
		{"latest", "eur", "", "", "gbp,usd", "https://api.exchangeratesapi.io/latest?symbols=GBP,USD&base=EUR"},
		{"history", "", "2019-01-01", "2019-01-20", "", "https://api.exchangeratesapi.io/history?start_at=2019-01-01&end_at=2019-01-20"},
		{"history", "gbp", "2019-01-01", "2019-01-20", "", "https://api.exchangeratesapi.io/history?start_at=2019-01-01&end_at=2019-01-20&base=GBP"},
		{"history", "", "2019-01-01", "2019-01-20", "gbp,usd", "https://api.exchangeratesapi.io/history?start_at=2019-01-01&end_at=2019-01-20&symbols=GBP,USD"},
		{"history", "EUR", "2019-01-01", "2019-01-20", "gbp,usd", "https://api.exchangeratesapi.io/history?start_at=2019-01-01&end_at=2019-01-20&symbols=GBP,USD&base=EUR"},
	}

	for _, test := range testcases {
		if item, _ := generateRequest(test.command, test.baseFlag, test.startFlag, test.endFlag, test.currency); item != test.expected {
			t.Error("Expected: ", test.expected, "But got: ", item)
		}
	}
}
