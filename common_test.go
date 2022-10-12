package hakucho

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

type testConfig struct {
	GrantUserEmail string `json:"grant-user-email"`
}

func loadTestConfig(t *testing.T) testConfig {
	t.Helper()

	cfgData, err := ioutil.ReadFile("testconfig.json")
	if err != nil {
		t.Log("no test config")
		return testConfig{}
	}

	var cfg testConfig

	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		t.Fatalf("Failed to unmarshal json: %s", err)
		return testConfig{}
	}

	return cfg
}
