# Rates

* Go program that queries https://exchangeratesapi.io/ to get historical and current currency rates

## Requirements

* Internet connection
* go 1.12

## Usage

### Latest Rates

```

./rates latest

// Changing base currency

./rates --base=GBP latest

// Returning only specified currencies

./rates  --currency=GBP,USD,EUR latest

```

### Historical Rates

```
./rates --start=yyyy-mm-dd --end=yyyy-mm-dd history
./rates  --start 2012-12-12 --end=2012-12-15 history

//changing base currency

./rates --start 2012-12-12 --end=2012-12-15 -base=USD history

// Returning specified currencies

./rates --start 2012-12-12 --end=2012-12-15 -currency=USD,EUR history
```

### List of Currencies

Please note that some currencies may return 0 if there is no current data provided for them. 

* IDR 
* DKK 
* INR 
* HRK 
* KRW 
* RUB 
* ZAR 
* HUF 
* MXN 
* ISK 
* CNY 
* USD 
* TRY 
* CZK 
* ILS 
* JPY 
* AUD 
* MYR 
* BRL 
* RON 
* PHP 
* CHF 
* SGD 
* BGN 
* NZD 
* THB 
* NOK 
* GBP 
* PLN 
* SEK 
* CAD 
* HKD 


## Building

```
go build .
```

## Testing 

```
go test .
```
