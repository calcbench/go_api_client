package calcbench_test

import (
	"os"
	"testing"

	"github.com/calcbench/go_api_client"
)

func TestCalcBench(t *testing.T) {
	c, err := calcbench.New()
	if err != nil {
		t.Errorf("error instantiating calcbench: %+v", err)
	}

	var authenticated bool
	if authenticated, err = c.Login(os.Getenv("CALCBENCH_USERNAME"), os.Getenv("CALCBENCH_PASSWORD")); err != nil {
		t.Errorf("error logging: %+v", err)
	}

	if !authenticated {
		t.Fatalf("error in credentials: %+v", err)
	}

	Query := calcbench.StandardizedDataRequest{
		CompaniesParameters: calcbench.CompaniesParameters{
			CompanyIdentifiers: []string{"MSFT"},
		},
		PeriodParameters: calcbench.PeriodParameters{
			Year:   2016,
			Period: 0,
		},
		PageParameters: calcbench.StandardizedDataPageParameters{
			Metrics: []string{"revenue"},
		},
	}

	response, err := c.StandardizedData(Query)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}
