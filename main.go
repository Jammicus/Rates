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
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
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

	for k, v := range r.Rates {
		fmt.Printf("Currency %s = %f\n", k, v)
	}

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
