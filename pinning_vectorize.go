package pinningcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type VectorizeFileCriteria struct {
	CID string
}

type VectorizeFileResult struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

type DeleteVectorizedFileCriteria struct {
	CID string
}

type DeleteVectorizedFileResult struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

type QueryVectorsCriteria struct {
	Text string
	TopK uint
}

type queryVectorsRequest struct {
	Text string `json:"text"`
	TopK uint   `json:"limit,omitempty"`
}

type QueryVectorsResult struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

func (pinning Pinning) vectorizeFile(ids PinningIdentifiers, criteria VectorizeFileCriteria) (VectorizeFileResult, []error) {
	var result VectorizeFileResult

	query := fmt.Sprintf("%s/v2/vectorize/files/%s", baseUrl, criteria.CID)
	log.Debug("POST ", query)

	req, err := http.NewRequest("POST", query, bytes.NewBuffer([]byte("{}")))
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
		return result, []error{fmt.Errorf("failed to call Pinning server with POST at %s", query),
			fmt.Errorf("status code: %s", resp.Status),
			fmt.Errorf("response body: %s", string(body)),
		}
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, []error{fmt.Errorf("failed to parse JSON response body while calling Pinning server with POST at %s", query), err}
	}

	return result, nil
}

func (pinning Pinning) deleteVectorizedFile(ids PinningIdentifiers, criteria DeleteVectorizedFileCriteria) (DeleteVectorizedFileResult, []error) {
	var result DeleteVectorizedFileResult

	query := fmt.Sprintf("%s/v2/vectorize/files/%s", baseUrl, criteria.CID)
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
		return result, []error{fmt.Errorf("failed to call Pinning server with DELETE at %s", query),
			fmt.Errorf("status code: %s", resp.Status),
			fmt.Errorf("response body: %s", string(body)),
		}
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, []error{fmt.Errorf("failed to parse JSON response body while calling Pinning server with DELETE at %s", query), err}
	}

	return result, nil
}

func (pinning Pinning) queryVectors(ids PinningIdentifiers, criteria QueryVectorsCriteria) (QueryVectorsResult, []error) {
	var result QueryVectorsResult

	request := queryVectorsRequest{
		Text: criteria.Text,
		TopK: criteria.TopK,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return result, []error{fmt.Errorf("failed to marshal query vectors request"), err}
	}

	query := fmt.Sprintf("%s/v2/vectorize/query", baseUrl)
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
		return result, []error{fmt.Errorf("failed to call Pinning server with POST at %s", query),
			fmt.Errorf("status code: %s", resp.Status),
			fmt.Errorf("response body: %s", string(body)),
		}
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, []error{fmt.Errorf("failed to parse JSON response body while calling Pinning server with POST at %s", query), err}
	}

	return result, nil
}
