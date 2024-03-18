package pinningcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type pinByHashRequest struct {
	HashToPin string    `json:"hash_to_pin"`
	Metadata  *Metadata `json:"metadata,omitempty"`
}

type pinByHashResult struct {
	Status string `json:"status"`
	Data   Pin    `json:"data"`
}

type PinByHashCriteria struct {
	HashToPin string
	Name      string
	KeyValues datatypes.JSON
}

func (criteria PinByHashCriteria) String() string {
	return fmt.Sprintf("HashToPin: %s", criteria.HashToPin)
}

func (pinning Pinning) buildPinByHashRequest(criteria PinByHashCriteria) pinByHashRequest {
	request := pinByHashRequest{
		HashToPin: criteria.HashToPin,
		Metadata: &Metadata{
			Name:      criteria.Name,
			KeyValues: criteria.KeyValues,
		},
	}

	return request
}

func (pinning Pinning) pinByHash(ids PinningIdentifiers, criteria PinByHashCriteria) (Pin, []error) {

	var result Pin

	request := pinning.buildPinByHashRequest(criteria)
	requestBody, err := json.Marshal(request)
	if err != nil {
		return result, []error{fmt.Errorf("failed to marshal criteria to JSON: %s", criteria), err}
	}

	query := fmt.Sprintf("%s/pinning/pinByHash/", baseUrl)

	log.Debug("POST ", query)
	log.Trace("Request body: ", string(requestBody))

	req, err := http.NewRequest("POST", query, bytes.NewBuffer(requestBody))
	if err != nil {
		return result, []error{fmt.Errorf("failed to prepare request for Pinning server with POST at %s", query), err}
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("pinning_api_key", ids.ApiKey)
	req.Header.Add("pinning_secret_key", ids.SecretKey)

	resp, err := httpClient.Do(req)

	if err != nil {
		return result, []error{fmt.Errorf("failed to call Pinning server with POST at %s", query), err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, []error{fmt.Errorf("failed to read body while calling Pinning server with POST at %s", query), err}
	}

	log.Trace("Status code: ", resp.Status)
	log.Trace("Response body: ", string(body))

	if resp.StatusCode != 200 {
		return result, []error{fmt.Errorf("frror while calling Pinning server with POST at %s", query),
			fmt.Errorf("status code: %s", resp.Status),
			fmt.Errorf("response body: %s", string(body)),
		}
	}

	var response pinByHashResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return result, []error{fmt.Errorf("failed to parse JSON response body while calling Pinning server with POST at %s", query), err}
	}

	result = Pin{
		ID:         response.Data.ID,
		CID:        response.Data.CID,
		UserID:     response.Data.UserID,
		Size:       response.Data.Size,
		DatePinned: response.Data.DatePinned,
		Metadata:   response.Data.Metadata,
	}

	return result, nil
}
