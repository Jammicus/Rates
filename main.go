package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var Api = "https://api.exchangeratesapi.io/"
var BaseFlag, StartFlag, EndFlag, CurrencyFlag *string

type Response struct {
	Base  string `json:"base"`
	Rates Rates  `json:"rates"`
	Date  string `json:"date"`
}

type Rates struct {
	USD float64 `json:"USD"`
	JPY float64 `json:"JPY"`
	BGN float64 `json:"BGN"`
	CZK float64 `json:"CZK"`
	DKK float64 `json:"DKK"`
	GPB float64 `json:"GPB"`
	HUF float64 `json:"HUF"`
	PLN float64 `json:"PLN"`
	RON float64 `json:"RON"`
	SEK float64 `json:"SEK"`
	CHF float64 `json:"CHF"`
	ISK float64 `json:"ISK"`
	NOK float64 `json:"NOK"`
	HRK float64 `json:"HRK"`
	RUB float64 `json:"RUB"`
	TRY float64 `json:"TRY"`
	BRL float64 `json:"BRL"`
	CAD float64 `json:"CAD"`
	CNY float64 `json:"CNY"`
	HKD float64 `json:"HKD"`
	IDR float64 `json:"IDR"`
	ILS float64 `json:"ILS"`
	INR float64 `json:"INR"`
	MXN float64 `json:"MXN"`
	MYR float64 `json:"MYR"`
	NZD float64 `json:"NZD"`
	PHP float64 `json:"PHP"`
	SGD float64 `json:"SGD"`
	THB float64 `json:"THB"`
	ZAR float64 `json:"ZAR"`
}

func main() {
	setLogging()
	req, err := generateRequest(flag.Arg(0), *BaseFlag, *StartFlag, *EndFlag, *CurrencyFlag)
	if err != nil {
		log.Fatal(err)
	}

	httpReq, err := request(req)
	output, err := parse(*httpReq)
	if err != nil {
		log.Fatal(err)
	}

	var r Response
	json.Unmarshal(output, &r)

	printResponce(r)
}

func printResponce(r Response) {
	fmt.Println("Date:", r.Date)
	fmt.Println("Base Currency:", r.Base)
	fmt.Println("Rates:")
	fmt.Printf("%+v\n", r.Rates)

}

func setLogging() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

func usage() {
	log.Errorln("Invalid flag, please use one of the following:")
	flag.PrintDefaults()
	os.Exit(2)
}

// Has to be flags then command, due to rules of the flag pakage.
func init() {
	BaseFlag = flag.String("base", "", "Specifies the base currency to use")
	StartFlag = flag.String("start", "", "Specifies the start date to use for a time series")
	EndFlag = flag.String("end", "", "Specifies the end date to use for a time series")
	CurrencyFlag = flag.String("currency", "", "Specifies a comma seperated list of currencies to be used")
	flag.Usage = usage
	flag.Parse()
}

func generateRequest(command, base, start, end, currency string) (string, error) {

	if command != "latest" && command != "history" {
		return "", errors.New("Invalid command " + command)
	}

	if start != "" && end == "" {
		return "", errors.New("Please provide the end flag when doing a time query")
	}
	if end != "" && start == "" {
		return "", errors.New("Please provide the start flag when doing a time query")

	}

	if command == "latest" && base == "" && currency == "" {
		return Api + command, nil
	}

	if command == "latest" && base != "" && currency == "" {
		return Api + command + "?" + "base=" + strings.ToUpper(base), nil
	}

	if command == "latest" && base == "" && currency != "" {
		return Api + command + "?" + "symbols=" + strings.ToUpper(currency), nil
	}

	if command == "latest" && base != "" && currency != "" {
		return Api + command + "?" + "symbols=" + strings.ToUpper(currency) + "&base=" + strings.ToUpper(base), nil
	}

	if command == "history" && start != "" && end != "" && base == "" && currency == "" {
		return Api + command + "?" + "start_at=" + start + "&end_at=" + end, nil
	}

	if command == "history" && start != "" && end != "" && base != "" && currency == "" {
		return Api + command + "?" + "start_at=" + start + "&end_at=" + end + "&" + "base=" + strings.ToUpper(base), nil
	}

	if command == "history" && start != "" && end != "" && base == "" && currency != "" {
		return Api + command + "?" + "start_at=" + start + "&end_at=" + end + "&symbols=" + strings.ToUpper(currency), nil
	}

	if command == "history" && start != "" && end != "" && base != "" && currency != "" {
		return Api + command + "?" + "start_at=" + start + "&end_at=" + end + "&symbols=" + strings.ToUpper(currency) + "&base=" + strings.ToUpper(base), nil
	}

	return "", errors.New("Unexpected request. Please try again")
}

func request(req string) (*http.Response, error) {
	resp, err := http.Get(req)
	return resp, err
}

func parse(resp http.Response) ([]byte, error) {
	var responseObject Response

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	errr := json.Unmarshal(responseData, &responseObject)
	return responseData, errr
}
