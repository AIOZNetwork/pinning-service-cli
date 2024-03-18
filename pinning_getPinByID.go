package pinningcli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type pinResponse struct {
	Data   Pin    `json:"data"`
	Status string `json:"status"`
}
type GetPinByIDCriteria struct {
	ID string
}

func (pinning Pinning) getPinByID(ids PinningIdentifiers, criteria GetPinByIDCriteria) (Pin, []error) {

	query := fmt.Sprintf("%s/pinning/%s", baseUrl, criteria.ID)
	log.Debug("GET ", query)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return Pin{}, []error{fmt.Errorf("failed to call Pinning server with GET at %s", query), err}
	}

	req.Header.Add("pinning_api_key", ids.ApiKey)
	req.Header.Add("pinning_secret_key", ids.SecretKey)

	resp, err := httpClient.Do(req)

	if err != nil {
		return Pin{}, []error{fmt.Errorf("failed to call Pinning server with GET at %s", query), err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pin{}, []error{fmt.Errorf("failed to read body while calling Pinning server with GET at %s", query), err}
	}

	log.Trace("Status code: ", resp.Status)
	log.Trace("Response body: ", string(body))

	if resp.StatusCode != 200 {
		return Pin{}, []error{fmt.Errorf("failed to call Pinning server with GET at %s", query),
			fmt.Errorf("response body: %s", string(body)),
		}
	}

	var response pinResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Pin{}, []error{fmt.Errorf("failed to parse JSON response body while calling Pinning server with GET at %s", query), err}
	}

	return response.Data, nil
}
