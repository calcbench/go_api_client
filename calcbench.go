package calcbench

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

//APIURL the base API URL for CalcBench API
const APIURL = "https://www.calcbench.com"

//CompaniesParameters specifies the companies to include
type CompaniesParameters struct {
	CompanyIdentifiers []string `json:"companyIdentifiers"`
	EntireUniverse     bool     `json:"entireUniverse,omitempty"`
}

//PeriodParameters specifies the periods to include
type PeriodParameters struct {
	Year       int32  `json:"year"`
	Period     int8   `json:"period"`
	EndYear    int    `json:"endYear,omitempty"`
	EndPeriod  int    `json:"endPeriod"`
	AllHistory bool   `json:"allHistory,omitempty"`
	UpdateDate string `json:"updateDate,omitempty"`
}

//StandardizedDataPageParameters specifies the metrics to search for
type StandardizedDataPageParameters struct {
	Metrics []string `json:"metrics"`
}

// StandardizedDataRequest for comparisons between companies
type StandardizedDataRequest struct {
	CompaniesParameters CompaniesParameters            `json:"companiesParameters"`
	PeriodParameters    PeriodParameters               `json:"periodParameters"`
	PageParameters      StandardizedDataPageParameters `json:"pageParameters"`
}

// StandardizedDataResponseObject result object
type StandardizedDataResponseObject struct {
	Ticker         string      `json:"ticker"`
	CalendarYear   int         `json:"calendar_year"`
	CalendarPeriod int         `json:"calendar_period"`
	Metric         string      `json:"metric"`
	Value          interface{} `json:"value"`
}

// StandardizedDataResponse is a slice of StandardizedDataResponseObject
type StandardizedDataResponse []StandardizedDataResponseObject

//CalcBench is the CalcBench API
type CalcBench struct {
	Client *http.Client
}

// New creates a new API instance
func New() (*CalcBench, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &CalcBench{
		Client: &http.Client{
			Jar: cookieJar,
		},
	}, nil
}

// Login authenticates the API using https://www.calcbench.com/api#authentication
func (c CalcBench) Login(email, password string) (bool, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("password", password)

	req, err := http.NewRequest("POST",
		APIURL+"/account/LogOnAjax",
		strings.NewReader(data.Encode()))
	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("invalid status code (%d)", resp.StatusCode)
	}

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result bool
	if err := json.Unmarshal(buffer, &result); err != nil {
		return false, err
	}

	return result, nil
}

//StandardizedData calls https://www.calcbench.com/api#normalizedData
func (c CalcBench) StandardizedData(Query StandardizedDataRequest) (StandardizedDataResponse, error) {
	b, err := json.Marshal(Query)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Query: %s\n", string(b))
	req, err := http.NewRequest("POST", APIURL+"/api/mappedData", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %s", string(contents))
	}

	var result StandardizedDataResponse
	if err := json.Unmarshal(contents, &result); err != nil {
		return nil, err
	}
	return result, nil
}
