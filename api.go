package nbrb

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"fmt"
	"time"
	"errors"
)

const (
	ENDPOINT string = "http://www.nbrb.by/API/ExRates"
)

type CurrencyDescription struct {
	Cur_ID            int
	Cur_ParentID      int
	Cur_Code          string
	Cur_Abbreviation  string
	Cur_Name          string
	Cur_Name_Bel      string
	Cur_Name_Eng      string
	Cur_QuotName      string
	Cur_QuotName_Bel  string
	Cur_QuotName_Eng  string
	Cur_NameMulti     string
	Cur_Name_BelMulti string
	Cur_Name_EngMulti string
	Cur_Scale         int
	Cur_Periodicity   int
	Cur_DateStart     string
	Cur_DateEnd       string
}

type Currency struct {
	Cur_ID           int
	Date             string
	Cur_Abbreviation string
	Cur_Scale        int
	Cur_Name         string
	Cur_OfficialRate float64
}

type Api struct {
	Currencies  []CurrencyDescription
	Periodicity int
}

func NewApi() *Api {
	return &Api{
		Periodicity:0,
	}
}

func (api *Api) call(url string) ([]byte, error) {
	cl := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	response, err := cl.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return body, errors.New("No Data")
	}

	return body, err
}

/*
  Get List of all currencies
 */
func (api *Api) GetCurrencies() (*Api, error) {

	url := api.currencyUrl()

	body, err := api.call(url)
	if err != nil {
		return api, err
	}

	err = json.Unmarshal(body, &api.Currencies)
	if err != nil {
		return api, err
	}
	return api, err
}

// Get URL for Currency
func (api *Api) currencyUrl() string {
	return fmt.Sprintf("%s/%s", ENDPOINT, "/Currencies")
}

/*
 Get specific currency form the API
 */
func (api *Api) GetCurrency(id int) (CurrencyDescription, error) {

	url := api.currencyUrl()
	url += "/" + strconv.Itoa(id)
	body, err := api.call(url)

	cd := CurrencyDescription{}
	if err != nil {
		return cd, err
	}

	err = json.Unmarshal(body, &cd)
	if err != nil {
		return cd, err
	}
	return cd, err
}

func (api *Api) GetRate(id int) (Currency, error) {
	url := api.rateUrl(id)
	rate := Currency{}

	body, err := api.call(url)
	if err != nil {
		return rate, err
	}

	err = json.Unmarshal(body, &rate)

	if err != nil {
		return rate, err
	}

	return rate, err
}

/*
   Get Currency Rate on specific Date
   The date format should be YYYY-M-D
 */
func (api *Api) GetRateOnDate(id int, date string) (Currency, error) {
	c := Currency{}

	date, err := api.validateDate(date)
	if err != nil {
		return c, err
	}
	url := fmt.Sprintf("%s?onDate=%s", api.rateUrl(id), date)
	resp, err := api.call(url)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(resp, &c)
	return c, err
}

// Date validation
func (api *Api) validateDate(i string) (string, error) {
	_, err := time.Parse("2006-1-2T15:04:05", i+"T00:00:00")
	return i, err
}

/*
  Get API Rate Url
 */
func (api *Api) rateUrl(i int) string {
	return fmt.Sprintf("%s/Rates/%d", ENDPOINT, i)
}
