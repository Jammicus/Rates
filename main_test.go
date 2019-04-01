package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
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

func newRequest(testDataFile string) *http.Response {
	fileContent, _ := ioutil.ReadFile(testDataFile)
	res := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(fileContent)),
	}
	return res
}

func TestLatestParseRequest(t *testing.T) {

	var latest = ResponseLatest{
		Base: "EUR",
		Date: "2019-03-29",
		Rates: map[string]float64{
			"AUD": 1.5821,
			"BGN": 1.9558,
			"BRL": 4.3865,
			"CAD": 1.5,
			"CHF": 1.1181,
			"CNY": 7.5397,
			"CZK": 25.802,
			"DKK": 7.4652,
			"GBP": 0.8583,
			"HKD": 8.8195,
			"HRK": 7.4338,
			"HUF": 321.05,
			"IDR": 15998.64,
			"ILS": 4.0764,
			"INR": 77.719,
			"ISK": 137.5,
			"JPY": 124.45,
			"KRW": 1276.46,
			"MXN": 21.691,
			"MYR": 4.5838,
			"NOK": 9.659,
			"NZD": 1.65,
			"PHP": 59.075,
			"PLN": 4.3006,
			"RON": 4.7608,
			"RUB": 72.8564,
			"SEK": 10.398,
			"SGD": 1.5214,
			"THB": 35.632,
			"TRY": 6.3446,
			"USD": 1.1235,
			"ZAR": 16.2642,
		},
	}
	var latestBase = ResponseLatest{
		Base: "EUR",
		Date: "2019-03-29",
		Rates: map[string]float64{
			"AUD": 1.582100,
			"BGN": 1.955800,
			"BRL": 4.386500,
			"CAD": 1.500000,
			"CHF": 1.118100,
			"CNY": 7.539700,
			"CZK": 25.802000,
			"DKK": 7.465200,
			"GBP": 0.858300,
			"HKD": 8.819500,
			"HRK": 7.433800,
			"HUF": 321.050000,
			"IDR": 15998.640000,
			"ILS": 4.076400,
			"INR": 77.719000,
			"ISK": 137.500000,
			"JPY": 124.450000,
			"KRW": 1276.460000,
			"MXN": 21.691000,
			"MYR": 4.583800,
			"NOK": 9.659000,
			"NZD": 1.650000,
			"PHP": 59.075000,
			"PLN": 4.300600,
			"RON": 4.760800,
			"RUB": 72.856400,
			"SEK": 10.398000,
			"SGD": 1.521400,
			"THB": 35.632000,
			"TRY": 6.344600,
			"USD": 1.123500,
			"ZAR": 16.264200,
		},
	}

	var testcases = []struct {
		filePath       string
		expectedStruct ResponseLatest
	}{
		{"testdata/latest.json", latest},
		{"testdata/latestBaseUSD.json", latestBase},
	}

	for _, test := range testcases {
		var l ResponseLatest
		r := newRequest(test.filePath)
		output, err := l.parseRequest(*r)
		if err != nil {
			log.Fatal(err)
		}
		x := output.rates()
		y := x.(map[string]float64)
		if !reflect.DeepEqual(y, test.expectedStruct.Rates) {
			t.Error("Expected: ", y, "\n", "But got: ", test.expectedStruct.Rates)
		}
	}
}

func TestHistoryParseRequest(t *testing.T) {

	var history = ResponseHistory{
		Base:  "EUR",
		Start: "2018-01-03",
		End:   "2018-01-01",
		Rates: map[string]Rates{
			"2018-01-03": Rates{
				BGN: 1.9558,
				NZD: 1.6942,
				ILS: 4.1588,
				RUB: 69.0962,
				CAD: 1.5047,
				USD: 1.2023,
				PHP: 59.988,
				CHF: 1.1736,
				ZAR: 14.8845,
				AUD: 1.5339,
				JPY: 134.97,
				TRY: 4.5303,
				HKD: 9.3985,
				MYR: 4.8272,
				THB: 39.11,
				HRK: 7.441,
				NOK: 9.744,
				IDR: 16176.95,
				DKK: 7.4442,
				CZK: 25.545,
				HUF: 309.29,
				GBP: 0.8864,
				MXN: 23.3835,
				KRW: 1281.39,
				SGD: 1.5988,
				BRL: 3.9236,
				PLN: 4.1652,
				INR: 76.3455,
				RON: 4.6355,
				CNY: 7.8168,
				SEK: 9.825,
			},
			"2018-01-02": Rates{
				BGN: 1.9558,
				NZD: 1.6955,
				ILS: 4.1693,
				RUB: 69.1176,
				CAD: 1.5128,
				USD: 1.2065,
				PHP: 60.132,
				CHF: 1.1718,
				ZAR: 14.9,
				AUD: 1.5413,
				JPY: 135.35,
				TRY: 4.534,
				HKD: 9.4283,
				MYR: 4.8495,
				THB: 39.115,
				HRK: 7.464,
				NOK: 9.7748,
				IDR: 16266.03,
				DKK: 7.4437,
				CZK: 25.494,
				HUF: 308.59,
				GBP: 0.88953,
				MXN: 23.5534,
				KRW: 1281.59,
				SGD: 1.6031,
				BRL: 3.9504,
				PLN: 4.1633,
				INR: 76.6005,
				RON: 4.6525,
				CNY: 7.8338,
				SEK: 9.8283,
			},
		},
	}

	var historyBase = ResponseHistory{
		Base:  "USD",
		Start: "2018-01-03",
		End:   "2018-01-01",
		Rates: map[string]Rates{
			"2018-01-03": Rates{
				BGN: 1.626715462,
				NZD: 1.409132496,
				ILS: 3.459036846,
				RUB: 57.470015803,
				CAD: 1.251517924,
				USD: 1,
				PHP: 49.8943691258,
				CHF: 0.9761290859,
				AUD: 1.2758047076,
				JPY: 112.2598353156,
				TRY: 3.7680279464,
				HKD: 7.8171005573,
				MYR: 4.014971305,
				HRK: 6.1889711387,
				CZK: 21.2467770107,
				IDR: 13455.0029110871,
				DKK: 6.191632704,
				NOK: 8.1044664393,
				HUF: 257.2486068369,
				GBP: 0.7372535973,
				MXN: 19.4489728021,
				THB: 32.5293188056,
				ZAR: 12.3800216252,
				BRL: 3.2634117941,
				SGD: 1.3297845796,
				PLN: 3.4643599767,
				INR: 63.4995425435,
				KRW: 1065.7822506862,
				RON: 3.8555269068,
				CNY: 6.5015387175,
				SEK: 8.1718373118,
				EUR: 0.8317391666,
			},
			"2018-01-02": Rates{
				BGN: 1.6210526316,
				NZD: 1.4053046001,
				ILS: 3.4556983009,
				RUB: 57.2876916701,
				CAD: 1.2538748446,
				USD: 1,
				PHP: 49.8400331538,
				CHF: 0.9712391214,
				AUD: 1.2774968918,
				JPY: 112.1840033154,
				TRY: 3.7579776212,
				HKD: 7.8145876502,
				MYR: 4.0194778284,
				HRK: 6.1864898467,
				CZK: 21.1305428927,
				IDR: 13481.9975134687,
				DKK: 6.1696643183,
				NOK: 8.1017820141,
				HUF: 255.772896809,
				GBP: 0.7372813925,
				MXN: 19.5220886863,
				THB: 32.4202237878,
				ZAR: 12.349772068,
				BRL: 3.2742644012,
				SGD: 1.3287194364,
				PLN: 3.4507252383,
				INR: 63.4898466639,
				KRW: 1062.23787816,
				RON: 3.8561956071,
				CNY: 6.4929962702,
				SEK: 8.1461251554,
				EUR: 0.828843763,
			},
		},
	}

	var testcases = []struct {
		filePath       string
		expectedStruct ResponseHistory
	}{
		{"testdata/historicalRange.json", history},
		{"testdata/historicalRangeBaseUSD.json", historyBase},
	}

	for _, test := range testcases {
		var h ResponseHistory
		r := newRequest(test.filePath)
		output, err := h.parseRequest(*r)
		if err != nil {
			log.Fatal(err)
		}
		x := output.rates()
		y := x.(map[string]Rates)
		if !reflect.DeepEqual(y, test.expectedStruct.Rates) {
			t.Error("Expected: ", y, "\n", "But got: ", test.expectedStruct.Rates)
		}
	}
}
