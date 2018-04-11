package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	data := url.Values{}
	data.Set("email", os.Getenv("CALCBENCH_USERNAME"))
	data.Set("strng", os.Getenv("CALCBENCH_PASSWORD"))
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	r, _ := http.NewRequest("POST", "https://www.calcbench.com/account/LogOnAjax", strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	}
	periodParams := PeriodParams{2016, 0}
	companiesParams := CompaniesParameters{[]string{"msft"}}
	pageParams := StandardizedDataPageParams{[]string{"revision_count"}}
	APIqueryParams := StandardizedDataParams{pageParams, periodParams, companiesParams}
	b, err := json.Marshal(APIqueryParams)
	if err != nil {
		fmt.Printf("%s", err)
	}
	req, err := http.NewRequest("POST", "https://www.calcbench.com/api/mappedData", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	}
}

type PeriodParams struct {
	Year   int32
	Period int8
}

type StandardizedDataPageParams struct {
	Metrics []string
}

type CompaniesParameters struct {
	CompanyIdentifiers []string
}

type StandardizedDataParams struct {
	PageParameters      StandardizedDataPageParams
	PeriodParameters    PeriodParams
	CompaniesParameters CompaniesParameters
}
