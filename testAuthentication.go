package pinningcli

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type TestAuthResult struct {
	Success bool `json:"success"`
}

func (pinning Pinning) testAuthentication(ids PinningIdentifiers) (TestAuthResult, []error) {

	var result = TestAuthResult{Success: false}

	query := fmt.Sprintf("%s/apiKeys/testAuthentication/", baseUrl)
	log.Debug("GET ", query)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return result, []error{fmt.Errorf("failed to prepare request for Pinning server with GET at %s", query), err}
	}

	req.Header.Add("pinning_api_key", ids.ApiKey)
	req.Header.Add("pinning_secret_key", ids.SecretKey)

	resp, err := httpClient.Do(req)

	if err != nil {
		return result, []error{fmt.Errorf("failed to call Pinning server with GET at %s", query), err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, []error{fmt.Errorf("failed to read body while calling Pinning server with GET at %s", query), err}
	}

	log.Trace("Status code: ", resp.Status)
	log.Trace("Response body: ", string(body))

	if resp.StatusCode != 200 {
		return result, nil
	}

	return TestAuthResult{Success: true}, nil
}
