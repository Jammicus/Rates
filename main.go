package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type Response interface {
	base() string
	date() map[string]string
	parseRequest(http.Response) (Response, error)
	printInfo()
	rates() interface{}
}

var Api = "https://api.exchangeratesapi.io/"
var BaseFlag, StartFlag, EndFlag, CurrencyFlag *string

type ResponseLatest struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

type ResponseHistory struct {
	Base  string           `json:"base"`
	Rates map[string]Rates `json:"rates"`
	End   string           `json:"end_at"`
	Start string           `json:"start_at"`
}

type Rates struct {
	AUD float64 `json:"AUD"`
	BGN float64 `json:"BGN"`
	BRL float64 `json:"BRL"`
	CAD float64 `json:"CAD"`
	CHF float64 `json:"CHF"`
	CNY float64 `json:"CNY"`
	CZK float64 `json:"CZK"`
	DKK float64 `json:"DKK"`
	EUR float64 `json:"EUR"`
	GBP float64 `json:"GBP"`
	HKD float64 `json:"HKD"`
	HRK float64 `json:"HRK"`
	HUF float64 `json:"HUF"`
	IDR float64 `json:"IDR"`
	ILS float64 `json:"ILS"`
	INR float64 `json:"INR"`
	ISK float64 `json:"ISK"`
	JPY float64 `json:"JPY"`
	KRW float64 `json:"KRW"`
	MXN float64 `json:"MXN"`
	MYR float64 `json:"MYR"`
	NOK float64 `json:"NOK"`
	NZD float64 `json:"NZD"`
	PHP float64 `json:"PHP"`
	PLN float64 `json:"PLN"`
	RON float64 `json:"RON"`
	RUB float64 `json:"RUB"`
	SEK float64 `json:"SEK"`
	SGD float64 `json:"SGD"`
	THB float64 `json:"THB"`
	TRY float64 `json:"TRY"`
	USD float64 `json:"USD"`
	ZAR float64 `json:"ZAR"`
}

func main() {
	setLogging()
	cmd := flag.Arg(0)
	r := determineResponseType(cmd)

	req, err := generateRequest(cmd, *BaseFlag, *StartFlag, *EndFlag, *CurrencyFlag)
	if err != nil {
		log.Fatal(err)
	}
	httpReq, err := sendRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	output, err := r.parseRequest(*httpReq)
	if err != nil {
		log.Fatal(err)
	}
	printResponce(output)

}

func determineResponseType(cmd string) Response {
	if cmd == "latest" {
		var r ResponseLatest
		return r
	}
	if cmd == "history" {
		var r ResponseHistory
		return r
	}
	return nil

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

func printResponce(r Response) {
	r.printInfo()
}

func setLogging() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

func usage() {
	log.Error("Invalid flag, please use one of the following:")
	flag.PrintDefaults()
	os.Exit(2)
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

func sendRequest(req string) (*http.Response, error) {
	resp, err := http.Get(req)
	return resp, err
}

func (r ResponseHistory) parseRequest(resp http.Response) (Response, error) {
	var responseHistory ResponseHistory
	var errr error

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseHistory, err
	}

	// Improve this.

	errr = json.Unmarshal(responseData, &responseHistory)
	return responseHistory, errr
}

func (r ResponseLatest) parseRequest(resp http.Response) (Response, error) {
	var responseLatest ResponseLatest
	var errr error

	responseData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return responseLatest, err
	}

	errr = json.Unmarshal(responseData, &responseLatest)
	return responseLatest, errr
}

// Think about how to remove reflection
func (r ResponseHistory) printInfo() {
	fmt.Println("Base Currency:", r.base())
	fmt.Println("Start Date:", r.date()["start"])
	fmt.Println("End Date:", r.date()["end"])

	for k, v := range r.rates().(map[string]Rates) {
		fmt.Println("")
		fmt.Println("Rates on the date of:", k)
		fmt.Println("")
		elem := reflect.ValueOf(&v).Elem()

		for i := 0; i < elem.NumField(); i++ {
			fmt.Printf("Currency %s = %v\n",
				elem.Type().Field(i).Name, elem.Field(i).Interface())
		}
	}
}

func (r ResponseLatest) printInfo() {
	fmt.Println("Base Currency:", r.base())
	fmt.Println("")
	fmt.Println("Date:", r.date()["date"])
	fmt.Println("")
	for k, v := range r.rates().(map[string]float64) {
		fmt.Printf("Currency %s = %f\n", k, v)
	}
}

func (r ResponseLatest) base() string {
	return r.Base
}

func (r ResponseLatest) rates() interface{} {
	return r.Rates
}

func (r ResponseLatest) date() map[string]string {
	m := make(map[string]string)
	m["date"] = r.Date
	return m
}

func (r ResponseHistory) base() string {
	return r.Base
}

func (r ResponseHistory) rates() interface{} {
	return r.Rates
}

func (r ResponseHistory) date() map[string]string {
	m := make(map[string]string)
	m["start"] = r.Start
	m["end"] = r.End
	return m
}
