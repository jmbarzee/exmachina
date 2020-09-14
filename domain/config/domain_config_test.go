package config

import (
	"testing"
	"time"
)

func TestDomainConfigCheck(t *testing.T) {
	type DommainConfigCheckTest struct {
		Config        DomainConfig
		ExpectedError error
	}

	cases := []DommainConfigCheckTest{
		{
			Config: DomainConfig{
				UUID: "DMAIN_0000000000000000",
				Traits: []string{
					"testTrait",
				},
				LogFilePath:    "/path/to/log/file.log",
				Port:           8080,
				DialTimeout:    time.Second,
				ServiceCheck:   time.Second,
				IsolationCheck: time.Second,
			},
		},
	}

	for _, c := range cases {
		err := c.Config.check()
		if err != c.ExpectedError {
			// t.Fatal
		}
	}
}
