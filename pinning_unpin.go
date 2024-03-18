package pinningcli

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type UnpinCriteria struct {
	ID string
}

type UnpinResult struct {
	Success bool `json:"success"`
}

func (criteria UnpinCriteria) String() string {
	return fmt.Sprintf("ID: %s", criteria.ID)
}

func (result UnpinResult) String() string {
	return fmt.Sprintf("Success: %s",
		strconv.FormatBool(result.Success))
}

func (pinning Pinning) unpin(ids PinningIdentifiers, criteria UnpinCriteria) (UnpinResult, []error) {

	var result = UnpinResult{Success: false}

	query := fmt.Sprintf("%s/pinning/unpin/%s", baseUrl, criteria.ID)
	log.Debug("DELETE ", query)

	req, err := http.NewRequest("DELETE", query, nil)
	if err != nil {
		return result, []error{fmt.Errorf("failed to prepare request for Pinning server with DELETE at %s", query), err}
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("pinning_api_key", ids.ApiKey)
	req.Header.Add("pinning_secret_key", ids.SecretKey)

	resp, err := httpClient.Do(req)

	if err != nil {
		return result, []error{fmt.Errorf("failed to call Pinning server with DELETE at %s", query), err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, []error{fmt.Errorf("failed to read body while calling Pinning server with DELETE at %s", query), err}
	}

	log.Trace("Status code: ", resp.Status)
	log.Trace("Response body: ", string(body))

	if resp.StatusCode != 200 {
		return result, []error{fmt.Errorf("response body: %s", string(body))}
	}

	return UnpinResult{Success: true}, nil
}
